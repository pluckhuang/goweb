package service

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/ecodeclub/ekit/slice"
	followv1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/follow/v1"
	"github.com/pluckhuang/goweb/aweb/feed/domain"
	"github.com/pluckhuang/goweb/aweb/feed/repository"
	"golang.org/x/sync/errgroup"
)

type ArticleEventHandler struct {
	repo         repository.FeedEventRepo
	followClient followv1.FollowServiceClient
}

const (
	ArticleEventName = "article_event"
	threshold        = 4
)

func NewArticleEventHandler(repo repository.FeedEventRepo, client followv1.FollowServiceClient) Handler {
	return &ArticleEventHandler{
		repo:         repo,
		followClient: client,
	}
}

func (h *ArticleEventHandler) FindFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	// article 这边是要聚合的
	// 按时间 事件可能在 push event，可能在 pull event
	var eg errgroup.Group
	var lock sync.Mutex
	events := make([]domain.FeedEvent, 0, limit*2)
	eg.Go(func() error {
		// 查询发件箱
		resp, err := h.followClient.GetFollowee(ctx, &followv1.GetFolloweeRequest{Follower: uid, Limit: 200})
		if err != nil {
			return err
		}
		followeeIDs := slice.Map(resp.FollowRelations, func(idx int, src *followv1.FollowRelation) int64 {
			return src.Followee
		})
		evts, err := h.repo.FindPullEventsWithTyp(ctx, ArticleEventName, followeeIDs, timestamp, limit)
		fmt.Println("FindFeedEvents: pull events:", len(evts), "uid:", uid, "timestamp:", timestamp, "limit:", limit)
		if err != nil {
			return err
		}
		lock.Lock()
		events = append(events, evts...)
		lock.Unlock()
		return nil
	})

	eg.Go(func() error {
		evts, err := h.repo.FindPushEventsWithTyp(ctx, ArticleEventName, uid, timestamp, limit)
		if err != nil {
			return err
		}
		lock.Lock()
		events = append(events, evts...)
		lock.Unlock()
		return nil
	})

	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	// 排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].Ctime.UnixMilli() > events[j].Ctime.UnixMilli()
	})
	return events[:slice.Min[int]([]int{int(limit), len(events)})], nil
}

func (h *ArticleEventHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {
	uid, err := ext.Get("uid").AsInt64()
	if err != nil {
		return err
	}
	// 找到这个人的粉丝数量，判定是拉模型还是推模型
	resp, err := h.followClient.GetFollowStatics(ctx, &followv1.GetFollowStaticsRequest{UserId: uid})
	if err != nil {
		return err
	}

	// 大v
	if resp.FollowStatics.Followers > threshold {
		// 拉模型
		return h.repo.CreatePullEvent(ctx, domain.FeedEvent{Uid: uid,
			Type:  ArticleEventName,
			Ctime: time.Now(),
			Ext:   ext})
	} else {
		// 推模型，也就是写扩散
		// 先查询出来粉丝
		fresp, err := h.followClient.GetFollower(ctx, &followv1.GetFollowerRequest{Followee: uid})
		if err != nil {
			return err
		}
		events := slice.Map(fresp.FollowRelations, func(idx int, src *followv1.FollowRelation) domain.FeedEvent {
			return domain.FeedEvent{Uid: src.Follower, Ctime: time.Now(), Type: ArticleEventName, Ext: ext}
		})
		return h.repo.CreatePushEvents(ctx, events)
	}
}
