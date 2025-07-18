package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/pluckhuang/goweb/aweb/comment/domain"
	"github.com/pluckhuang/goweb/aweb/comment/repository/dao"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type CommentRepository interface {
	// FindByBiz 根据 ID 倒序查找
	// 并且会返回每个评论的三条直接回复
	FindByBiz(ctx context.Context, biz string,
		bizId, minID, limit int64) ([]domain.Comment, bool, error)
	// 删除评论，删除本评论何其子评论
	DeleteComment(ctx context.Context, comment domain.Comment) error
	// CreateComment 创建评论
	CreateComment(ctx context.Context, comment domain.Comment) error
	// 获取更多的一级评论对应的子评论
	GetMoreReplies(ctx context.Context, rid int64, id int64, limit int64) ([]domain.Comment, bool, error)
}

type CachedCommentRepo struct {
	dao dao.CommentDAO
	l   logger.LoggerV1
}

func NewCommentRepo(commentDAO dao.CommentDAO, l logger.LoggerV1) CommentRepository {
	return &CachedCommentRepo{
		dao: commentDAO,
		l:   l,
	}
}

func (c *CachedCommentRepo) FindByBiz(ctx context.Context, biz string,
	bizId, minID, limit int64) ([]domain.Comment, bool, error) {
	daoComments, hasMore, err := c.dao.FindByBiz(ctx, biz, bizId, minID, limit)
	if err != nil {
		return nil, false, err
	}
	res := make([]*domain.Comment, 0, len(daoComments))
	// 找子评论了，找三条
	var eg errgroup.Group
	downgrade := ctx.Value("downgrade") == "true"
	for _, dc := range daoComments {
		cm := c.toDomain(dc)
		cmPtr := &cm
		res = append(res, cmPtr)
		if downgrade {
			continue
		}
		eg.Go(func() error {
			subComments, err := c.dao.FindRepliesByPid(ctx, dc.Id, 0, 3)
			if err != nil {
				return err
			}
			cmPtr.Children = make([]domain.Comment, 0, len(subComments))
			for _, sc := range subComments {
				cmPtr.Children = append(cmPtr.Children, c.toDomain(sc))
			}
			return nil
		})
	}
	// Convert []*domain.Comment back to []domain.Comment before returning
	comments := make([]domain.Comment, 0, len(res))
	for _, cmPtr := range res {
		comments = append(comments, *cmPtr)
	}
	return comments, hasMore, eg.Wait()
}

func (c *CachedCommentRepo) DeleteComment(ctx context.Context, comment domain.Comment) error {
	return c.dao.Delete(ctx, dao.Comment{
		Id: comment.Id,
	})
}

func (c *CachedCommentRepo) CreateComment(ctx context.Context, comment domain.Comment) error {
	return c.dao.Insert(ctx, c.toEntity(comment))
}

func (c *CachedCommentRepo) GetMoreReplies(ctx context.Context, rid int64, maxID int64, limit int64) ([]domain.Comment, bool, error) {
	cs, hasMore, err := c.dao.FindRepliesByRid(ctx, rid, maxID, limit)
	if err != nil {
		return nil, false, err
	}
	res := make([]domain.Comment, 0, len(cs))
	for _, cm := range cs {
		res = append(res, c.toDomain(cm))
	}
	return res, hasMore, nil
}

func (c *CachedCommentRepo) toDomain(daoComment dao.Comment) domain.Comment {
	val := domain.Comment{
		Id: daoComment.Id,
		Commentator: domain.User{
			ID: daoComment.Uid,
		},
		Biz:     daoComment.Biz,
		BizID:   daoComment.BizID,
		Content: daoComment.Content,
		CTime:   time.UnixMilli(daoComment.Ctime),
		UTime:   time.UnixMilli(daoComment.Utime),
	}
	if daoComment.PID.Valid {
		val.ParentComment = &domain.Comment{
			Id: daoComment.PID.Int64,
		}
	}
	if daoComment.RootID.Valid {
		val.RootComment = &domain.Comment{
			Id: daoComment.RootID.Int64,
		}
	}
	return val
}

func (c *CachedCommentRepo) toEntity(domainComment domain.Comment) dao.Comment {
	daoComment := dao.Comment{
		Id:      domainComment.Id,
		Uid:     domainComment.Commentator.ID,
		Biz:     domainComment.Biz,
		BizID:   domainComment.BizID,
		Content: domainComment.Content,
	}
	if domainComment.RootComment != nil {
		daoComment.RootID = sql.NullInt64{
			Valid: true,
			Int64: domainComment.RootComment.Id,
		}
	}
	if domainComment.ParentComment != nil {
		daoComment.PID = sql.NullInt64{
			Valid: true,
			Int64: domainComment.ParentComment.Id,
		}
	}
	daoComment.Ctime = time.Now().UnixMilli()
	daoComment.Utime = time.Now().UnixMilli()
	return daoComment
}
