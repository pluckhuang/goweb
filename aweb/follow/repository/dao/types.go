package dao

import "context"

// 存储用户的关注数据
type FollowRelation struct {
	ID int64 `gorm:"column:id;autoIncrement;primaryKey;"`

	// 要在这两个列上，创建一个联合唯一索引
	Follower int64 `gorm:"uniqueIndex:follower_followee"`
	Followee int64 `gorm:"uniqueIndex:follower_followee;index:idx_followee"`

	// 软删除策略
	Status uint8
	Ctime  int64
	Utime  int64
}

const (
	FollowRelationStatusUnknown uint8 = iota
	FollowRelationStatusActive
	FollowRelationStatusInactive
)

type FollowRelationDao interface {
	// CreateFollowRelation 创建联系人
	CreateFollowRelation(ctx context.Context, c FollowRelation) error
	// UpdateStatus 更新状态
	UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error
	// FollowerRelationList 获取某人的关注列表
	FollowerRelationList(ctx context.Context, follower, offset, limit int64) ([]FollowRelation, error)
	// FolloweeRelationList 获取某人的粉丝列表
	FolloweeRelationList(ctx context.Context, followee, offset, limit int64) ([]FollowRelation, error)
	// FollowRelationDetail 获取某个关注关系的详情
	FollowRelationDetail(ctx context.Context, follower int64, followee int64) (FollowRelation, error)
	// CntFollower 统计计算关注自己的人有多少
	CntFollower(ctx context.Context, uid int64) (int64, error)
	// CntFollowee 统计自己关注了多少人
	CntFollowee(ctx context.Context, uid int64) (int64, error)
}
