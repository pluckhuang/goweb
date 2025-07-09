package dao

import "context"

// FeedPullEventDAO 拉模型
type FeedPullEventDAO interface {
	// 创建拉取事件
	CreatePullEvent(ctx context.Context, event FeedPullEvent) error
	FindPullEventListWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error)
}

type FeedPushEventDAO interface {
	// CreatePushEvents 创建推送事件
	CreatePushEvents(ctx context.Context, events []FeedPushEvent) error
	FindPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp, limit int64) ([]FeedPushEvent, error)
}
