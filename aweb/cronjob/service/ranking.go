package service

import (
	"context"
	"math"
	"time"

	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	articlev1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/article/v1"
	interactivev1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/interactive/v1"
	"github.com/pluckhuang/goweb/aweb/article/domain"
	"github.com/pluckhuang/goweb/aweb/cronjob/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RankingService interface {
	// TopN 前 100 的
	TopN(ctx context.Context) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type BatchRankingService struct {
	// 用来取点赞数
	interactiveSvc interactivev1.InteractiveServiceClient

	// 用来查找文章
	artSvc articlev1.ArticleServiceClient

	batchSize int
	scoreFunc func(likeCnt int64, utime time.Time) float64
	n         int

	repo repository.RankingRepository
}

func NewBatchRankingService(interactiveSvc interactivev1.InteractiveServiceClient,
	artSvc articlev1.ArticleServiceClient, repo repository.RankingRepository) RankingService {
	return &BatchRankingService{
		interactiveSvc: interactiveSvc,
		artSvc:         artSvc,
		batchSize:      100,
		n:              100,
		repo:           repo,
		scoreFunc: func(likeCnt int64, utime time.Time) float64 {
			// 时间
			duration := time.Since(utime).Seconds()
			return float64(likeCnt-1) / math.Pow(duration+2, 1.5)
		},
	}
}

func (b *BatchRankingService) GetTopN(ctx context.Context) ([]domain.Article, error) {
	return b.repo.GetTopN(ctx)
}

func (b *BatchRankingService) TopN(ctx context.Context) error {
	arts, err := b.topN(ctx)
	if err != nil {
		return err
	}
	// 存到缓存里
	return b.repo.ReplaceTopN(ctx, arts)
}

func (b *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
	offset := 0
	start := time.Now()
	ddl := start.Add(-7 * 24 * time.Hour)

	type Score struct {
		score float64
		art   domain.Article
	}
	topN := queue.NewPriorityQueue[Score](b.n,
		func(src Score, dst Score) int {
			if src.score > dst.score {
				return 1
			} else if src.score == dst.score {
				return 0
			} else {
				return -1
			}
		})

	for {
		// 取数据
		artsResp, err := b.artSvc.ListPub(ctx, &articlev1.ListPubRequest{
			Start:  timestamppb.New(start),
			Offset: int32(offset),
			Limit:  int32(b.batchSize),
		})
		if err != nil {
			return nil, err
		}
		arts := artsResp.Articles
		ids := slice.Map(arts, func(idx int, art *articlev1.Article) int64 {
			return art.Id
		})
		interactiveResp, err := b.interactiveSvc.GetByIds(ctx, &interactivev1.GetByIdsRequest{
			Biz: "article", Ids: ids,
		})
		if err != nil {
			return nil, err
		}
		interactiveMap := interactiveResp.Intrs
		for _, art := range arts {
			intr := interactiveMap[art.Id]
			score := b.scoreFunc(intr.LikeCnt, art.Utime.AsTime())
			ele := Score{
				score: score,
				art: domain.Article{
					Id:      art.Id,
					Title:   art.Title,
					Content: art.Content,
					Author: domain.Author{
						Id: art.AuthorId,
					},
					Utime: art.Utime.AsTime(),
					Ctime: art.Ctime.AsTime(),
				},
			}
			err = topN.Enqueue(ele)
			if err == queue.ErrOutOfCapacity {
				// 这个也是满了
				// 拿出最小的元素
				minEle, _ := topN.Dequeue()
				if minEle.score < score {
					_ = topN.Enqueue(ele)
				} else {
					_ = topN.Enqueue(minEle)
				}
			}
		}
		offset = offset + len(arts)
		// 没有取够一批，我们就直接中断执行
		// 没有下一批了
		if len(arts) < b.batchSize ||
			// 这个是一个优化
			arts[len(arts)-1].Utime.AsTime().Before(ddl) {
			break
		}
	}

	// 这边 topN 里面就是最终结果
	res := make([]domain.Article, topN.Len())
	for i := topN.Len() - 1; i >= 0; i-- {
		ele, _ := topN.Dequeue()
		res[i] = ele.art
	}
	return res, nil
}
