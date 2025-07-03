package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GORMFollowRelationDAO struct {
	db *gorm.DB
}

func NewGORMFollowRelationDAO(db *gorm.DB) FollowRelationDao {
	return &GORMFollowRelationDAO{
		db: db,
	}
}

func (g *GORMFollowRelationDAO) CreateFollowRelation(ctx context.Context, f FollowRelation) error {
	// 我也要保持 insert or update 语义
	now := time.Now().UnixMilli()
	f.Ctime = now
	f.Utime = now
	f.Status = FollowRelationStatusActive
	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			// 这代表的是关注了-取消了-再关注了
			"status": FollowRelationStatusActive,
			"utime":  now,
		}),
	}).Create(&f).Error
}

func (g *GORMFollowRelationDAO) UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error {
	return g.db.WithContext(ctx).
		Where("follower = ? AND followee = ?", follower, followee).
		Updates(map[string]any{
			"status": status,
			"utime":  time.Now().UnixMilli(),
		}).Error
}

func (g *GORMFollowRelationDAO) CntFollower(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := g.db.WithContext(ctx).
		Select("count(follower)").
		// 如果要是没有额外索引，不用怀疑，全表扫描
		// 可以考虑在 followee 额外创建一个索引
		Where("followee = ? AND status = ?",
			uid, FollowRelationStatusActive).Count(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := g.db.WithContext(ctx).
		Select("count(followee)").
		// <follower, followee>
		Where("follower = ? AND status = ?",
			uid, FollowRelationStatusActive).Count(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) FollowRelationList(ctx context.Context,
	// 可以考虑只查询 follower id，于是命中了覆盖索引
	follower, offset, limit int64) ([]FollowRelation, error) {
	var res []FollowRelation
	err := g.db.WithContext(ctx).
		Where("follower = ? AND status = ?", follower, FollowRelationStatusActive).
		Offset(int(offset)).Limit(int(limit)).
		Find(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) FollowRelationDetail(ctx context.Context, follower int64, followee int64) (FollowRelation, error) {
	var res FollowRelation
	err := g.db.WithContext(ctx).Where("follower = ? AND followee = ? AND status = ?",
		follower, followee, FollowRelationStatusActive).First(&res).Error
	return res, err
}
