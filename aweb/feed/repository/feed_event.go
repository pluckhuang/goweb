package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pluckhuang/goweb/aweb/feed/domain"
	"github.com/pluckhuang/goweb/aweb/feed/repository/dao"
)

type FeedEventRepo interface {
	// CreatePushEvents 批量推事件
	CreatePushEvents(ctx context.Context, events []domain.FeedEvent) error
	// CreatePullEvent 创建拉事件
	CreatePullEvent(ctx context.Context, event domain.FeedEvent) error
	// FindPullEventsWithTyp 获取某个类型的拉事件，
	FindPullEventsWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]domain.FeedEvent, error)
	// FindPushEvents 获取某个类型的推事件，也就
	FindPushEventsWithTyp(ctx context.Context, typ string, uid, timestamp, limit int64) ([]domain.FeedEvent, error)
}

type feedEventRepo struct {
	pullDao dao.FeedPullEventDAO
	pushDao dao.FeedPushEventDAO
}

func NewFeedEventRepo(pullDao dao.FeedPullEventDAO, pushDao dao.FeedPushEventDAO) FeedEventRepo {
	return &feedEventRepo{
		pullDao: pullDao,
		pushDao: pushDao,
	}
}

func (f *feedEventRepo) CreatePushEvents(ctx context.Context, events []domain.FeedEvent) error {
	pushEvents := make([]dao.FeedPushEvent, 0, len(events))
	for _, e := range events {
		pushEvents = append(pushEvents, convertToPushEventDao(e))
	}
	return f.pushDao.CreatePushEvents(ctx, pushEvents)
}

func convertToPushEventDao(event domain.FeedEvent) dao.FeedPushEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPushEvent{
		Id:      event.ID,
		UID:     event.Uid,
		Type:    event.Type,
		Content: string(val),
		Ctime:   event.Ctime.Unix(),
	}
}

func (f *feedEventRepo) CreatePullEvent(ctx context.Context, event domain.FeedEvent) error {
	return f.pullDao.CreatePullEvent(ctx, convertToPullEventDao(event))
}

func convertToPullEventDao(event domain.FeedEvent) dao.FeedPullEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPullEvent{
		Id:      event.ID,
		UID:     event.Uid,
		Type:    event.Type,
		Content: string(val),
		Ctime:   event.Ctime.Unix(),
	}
}

func (f *feedEventRepo) FindPullEventsWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pullDao.FindPullEventListWithTyp(ctx, typ, uids, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPullEventDomain(e))
	}
	return ans, nil
}

func convertToPullEventDomain(event dao.FeedPullEvent) domain.FeedEvent {
	var ext map[string]string
	_ = json.Unmarshal([]byte(event.Content), &ext)
	return domain.FeedEvent{
		ID:    event.Id,
		Uid:   event.UID,
		Type:  event.Type,
		Ctime: time.Unix(event.Ctime, 0),
		Ext:   ext,
	}
}

func (f *feedEventRepo) FindPushEventsWithTyp(ctx context.Context, typ string, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pushDao.FindPushEventsWithTyp(ctx, typ, uid, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPushEventDomain(e))
	}
	return ans, nil
}

func convertToPushEventDomain(event dao.FeedPushEvent) domain.FeedEvent {
	var ext map[string]string
	_ = json.Unmarshal([]byte(event.Content), &ext)
	return domain.FeedEvent{
		ID:    event.Id,
		Uid:   event.UID,
		Type:  event.Type,
		Ctime: time.Unix(event.Ctime, 0),
		Ext:   ext,
	}
}
