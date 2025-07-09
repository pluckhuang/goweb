package dao

import (
	"context"

	"gorm.io/gorm"
)

// FeedPushEvent 对应的是收件箱
type FeedPushEvent struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 收件人
	UID  int64 `gorm:"index;column:uid"`
	Type string
	// 扩展字段，不同的事件类型，有不同的解析方式
	Content string
	Ctime   int64
}

type feedPushEventDAO struct {
	db *gorm.DB
}

func NewFeedPushEventDAO(db *gorm.DB) FeedPushEventDAO {
	return &feedPushEventDAO{
		db: db,
	}
}

func (f *feedPushEventDAO) FindPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp, limit int64) ([]FeedPushEvent, error) {
	var events []FeedPushEvent
	err := f.db.WithContext(ctx).
		Where("uid = ?", uid).
		Where("ctime < ?", timestamp).
		Where("type = ?", typ).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

func (f *feedPushEventDAO) CreatePushEvents(ctx context.Context, events []FeedPushEvent) error {
	return f.db.WithContext(ctx).Create(events).Error
}
