package service

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/ecodeclub/ekit/slice"
	followv1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/follow/v1"
	"github.com/pluckhuang/goweb/aweb/feed/domain"
	"github.com/pluckhuang/goweb/aweb/feed/repository"
	"golang.org/x/sync/errgroup"
)

type feedService struct {
	repo repository.FeedEventRepo
	// 对应的 string 就是 type
	handlerMap   map[string]Handler
	followClient followv1.FollowServiceClient
}

func NewFeedService(repo repository.FeedEventRepo, handlerMap map[string]Handler) FeedService {
	return &feedService{
		repo:       repo,
		handlerMap: handlerMap,
	}
}

func (f *feedService) RegisterService(typ string, handler Handler) {
	f.handlerMap[typ] = handler
}

func (f *feedService) CreateFeedEvent(ctx context.Context, feed domain.FeedEvent) error {
	handler, ok := f.handlerMap[feed.Type]
	if !ok {
		// 说明 type 不对
		return fmt.Errorf("未能找到对应的 Handler %s", feed.Type)
	}
	return handler.CreateFeedEvent(ctx, feed.Ext)
}

func (f *feedService) GetFeedEventList(ctx context.Context, uid int64, timestamp, limit int64) ([]domain.FeedEvent, error) {
	var eg errgroup.Group
	var lock sync.Mutex
	events := make([]domain.FeedEvent, 0, limit*int64(len(f.handlerMap)))
	for _, handler := range f.handlerMap {
		h := handler
		eg.Go(func() error {
			evts, err := h.FindFeedEvents(ctx, uid, timestamp, limit)
			if err != nil {
				return err
			}
			lock.Lock()
			events = append(events, evts...)
			lock.Unlock()
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	// 查询所有的数据，现在要排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].Ctime.UnixMilli() > events[j].Ctime.UnixMilli()
	})
	return events[:slice.Min[int]([]int{int(limit), len(events)})], nil
}
