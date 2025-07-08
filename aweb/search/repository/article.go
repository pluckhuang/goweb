package repository

import (
	"context"

	"github.com/ecodeclub/ekit/slice"
	"github.com/pluckhuang/goweb/aweb/search/domain"
	"github.com/pluckhuang/goweb/aweb/search/repository/dao"
	"golang.org/x/sync/errgroup"
)

type articleRepository struct {
	dao      dao.ArticleDAO
	collects dao.CollectDAO
	likes    dao.LikeDAO
}

func NewArticleRepository(d dao.ArticleDAO, collectDao dao.CollectDAO, like dao.LikeDAO) ArticleRepository {
	return &articleRepository{
		dao:      d,
		collects: collectDao,
		likes:    like,
	}
}

func (a *articleRepository) SearchArticle(ctx context.Context,
	uid int64,
	keywords []string) ([]domain.Article, error) {
	var eg errgroup.Group
	var collectArtIds, likeArtIds []int64
	var err error
	eg.Go(func() error {
		likeArtIds, err = a.likes.Search(ctx, uid, "article")
		return err
	})
	eg.Go(func() error {
		collectArtIds, err = a.collects.Search(ctx, uid, "article")
		return err
	})
	if err = eg.Wait(); err != nil {
		return nil, err
	}
	arts, err := a.dao.Search(ctx, dao.SearchReq{
		LikeIds:    likeArtIds,
		CollectIds: collectArtIds,
	}, keywords)
	if err != nil {
		return nil, err
	}
	return slice.Map(arts, func(idx int, src dao.Article) domain.Article {
		return domain.Article{
			Id:      src.Id,
			Title:   src.Title,
			Status:  src.Status,
			Content: src.Content,
		}
	}), nil
}

func (a *articleRepository) InputArticle(ctx context.Context, msg domain.Article) error {
	return a.dao.InputArticle(ctx, dao.Article{
		Id:      msg.Id,
		Title:   msg.Title,
		Status:  msg.Status,
		Content: msg.Content,
	})
}
