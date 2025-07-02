package dao

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

// ErrDataNotFound 通用的数据没找到
var ErrDataNotFound = gorm.ErrRecordNotFound

type CommentDAO interface {
	Insert(ctx context.Context, u Comment) error
	// FindByBiz 只查找一级评论
	FindByBiz(ctx context.Context, biz string,
		bizId, minID, limit int64) ([]Comment, error)
	// 查找一级评论对应的子评论
	FindRepliesByPid(ctx context.Context, pid int64, offset, limit int) ([]Comment, error)
	// Delete 删除本节点和其对应的子节点
	Delete(ctx context.Context, u Comment) error
	// 根据根评论的id和当前评论的id查找对应的回复
	FindRepliesByRid(ctx context.Context, rid int64, id int64, limit int64) ([]Comment, error)
}

// Comment 把这个评论的表结构设计好
type Comment struct {
	Id int64 `gorm:"autoIncrement,primaryKey"`
	// 发表评论的人
	Uid int64
	// 被评价的东西
	Biz     string `gorm:"index:biz_type_id"`
	BizID   int64  `gorm:"index:biz_type_id"`
	Content string

	// 根评论
	// 如果这个字段是 NULL，它是根评论
	RootID sql.NullInt64 `gorm:"column:root_id;index"`

	// 这个是 NULL，也是根评论
	PID sql.NullInt64 `gorm:"column:pid;index"`

	// 定义外键, 用于级联删除子评论
	ParentComment *Comment `gorm:"ForeignKey:PID;AssociationForeignKey:ID;constraint:OnDelete:CASCADE"`

	Ctime int64
	Utime int64
}

type GORMCommentDAO struct {
	db *gorm.DB
}

func NewCommentDAO(db *gorm.DB) CommentDAO {
	return &GORMCommentDAO{
		db: db,
	}
}

func (c *GORMCommentDAO) Insert(ctx context.Context, u Comment) error {
	return c.db.
		WithContext(ctx).
		Create(&u).
		Error
}

func (c *GORMCommentDAO) FindByBiz(ctx context.Context, biz string,
	bizId, minID, limit int64) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND id < ? AND pid IS NULL", biz, bizId, minID).
		Limit(int(limit)).
		Order("id DESC").
		Find(&res).Error
	return res, err
}

// FindRepliesByPid 查找评论的直接评论
func (c *GORMCommentDAO) FindRepliesByPid(ctx context.Context,
	pid int64,
	offset,
	limit int) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).Where("pid = ?", pid).
		Order("id DESC").
		Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}
func (c *GORMCommentDAO) Delete(ctx context.Context, u Comment) error {
	return c.db.WithContext(ctx).Delete(&Comment{
		Id: u.Id,
	}).Error
}

func (c *GORMCommentDAO) FindRepliesByRid(ctx context.Context,
	rid int64, id int64, limit int64) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).
		Where("root_id = ? AND id > ?", rid, id).
		Order("id ASC").
		Limit(int(limit)).Find(&res).Error
	return res, err
}
