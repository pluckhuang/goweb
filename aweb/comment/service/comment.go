package service

import (
	"context"

	"github.com/pluckhuang/goweb/aweb/comment/domain"
	"github.com/pluckhuang/goweb/aweb/comment/repository"
)

type CommentService interface {
	// 获取一级评论 按照 ID 倒序排序
	GetCommentList(ctx context.Context, biz string, bizId, minID, limit int64) ([]domain.Comment, bool, error)
	// 删除评论，删除本评论何其子评论
	DeleteComment(ctx context.Context, id int64) error
	//  创建评论
	CreateComment(ctx context.Context, comment domain.Comment) error
	// 获取更多的一级评论对应的子评论
	GetMoreReplies(ctx context.Context, rid int64, maxID int64, limit int64) ([]domain.Comment, bool, error)
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentSvc(repo repository.CommentRepository) CommentService {
	return &commentService{
		repo: repo,
	}
}

func (c *commentService) CreateComment(ctx context.Context, comment domain.Comment) error {
	return c.repo.CreateComment(ctx, comment)
}

func (c *commentService) DeleteComment(ctx context.Context, id int64) error {
	return c.repo.DeleteComment(ctx, domain.Comment{
		Id: id,
	})
}

func (c *commentService) GetCommentList(ctx context.Context, biz string,
	bizId, minID, limit int64) ([]domain.Comment, bool, error) {
	list, hasMore, err := c.repo.FindByBiz(ctx, biz, bizId, minID, limit)
	if err != nil {
		return nil, false, err
	}
	return list, hasMore, err
}

func (c *commentService) GetMoreReplies(ctx context.Context,
	rid int64,
	maxID int64, limit int64) ([]domain.Comment, bool, error) {
	return c.repo.GetMoreReplies(ctx, rid, maxID, limit)
}
