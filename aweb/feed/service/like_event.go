package service

import (
	"context"
	"time"

	"github.com/pluckhuang/goweb/aweb/feed/domain"
	"github.com/pluckhuang/goweb/aweb/feed/repository"
)

const (
	LikeEventName = "like_event"
)

type LikeEventHandler struct {
	repo repository.FeedEventRepo
}

func NewLikeEventHandler(repo repository.FeedEventRepo) Handler {
	return &LikeEventHandler{
		repo: repo,
	}
}

func (l *LikeEventHandler) FindFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	return l.repo.FindPushEventsWithTyp(ctx, LikeEventName, uid, timestamp, limit)
}

// CreateFeedEvent 中的 ext 里面至少需要三个 id
// liked int64: 被点赞的人
// liker int64：点赞的人
// bizId int64: 被点赞的东西
// biz: string: 业务类型，比如说 article, comment, video 等等
func (l *LikeEventHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {

	uid, err := ext.Get("liked").AsInt64()
	if err != nil {
		return err
	}
	return l.repo.CreatePushEvents(ctx, []domain.FeedEvent{{
		Uid:   uid,
		Ext:   ext,
		Ctime: time.Now(),
		Type:  LikeEventName}})
}
