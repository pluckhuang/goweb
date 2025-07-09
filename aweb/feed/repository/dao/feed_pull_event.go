package dao

import (
	"context"

	"gorm.io/gorm"
)

// FeedPullEvent 对应的是发件箱
type FeedPullEvent struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 发件人
	UID int64 `gorm:"index;column:uid"`
	// 事件类型，比如 xxx_event
	Type string
	// 扩展字段，不同的事件类型，有不同的解析方式
	Content string
	Ctime   int64
}

type feedPullEventDAO struct {
	db *gorm.DB
}

func NewFeedPullEventDAO(db *gorm.DB) FeedPullEventDAO {
	return &feedPullEventDAO{
		db: db,
	}
}

func (f *feedPullEventDAO) FindPullEventListWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error) {
	var events []FeedPullEvent
	err := f.db.WithContext(ctx).
		Where("uid in ?", uids).
		Where("ctime < ?", timestamp).
		Where("type = ?", typ).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

func (f *feedPullEventDAO) CreatePullEvent(ctx context.Context, event FeedPullEvent) error {
	return f.db.WithContext(ctx).Create(&event).Error
}
