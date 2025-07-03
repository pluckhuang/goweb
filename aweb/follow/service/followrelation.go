package service

import (
	"context"

	"github.com/pluckhuang/goweb/aweb/follow/domain"
	"github.com/pluckhuang/goweb/aweb/follow/repository"
)

type FollowRelationService interface {
	Follow(ctx context.Context, follower, followee int64) error
	CancelFollow(ctx context.Context, follower, followee int64) error
	// 获得某个人的关注列表
	GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error)
	// 获得某个人的粉丝列表
	GetFollower(ctx context.Context, followee, offset, limit int64) ([]domain.FollowRelation, error)
	// 获取某个人是否存在关注另外一个人的信息
	FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error)
	// GetFollowStatics 获取关注和粉丝数量
	GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error)
}

type followRelationService struct {
	repo repository.FollowRepository
}

func NewFollowRelationService(repo repository.FollowRepository) FollowRelationService {
	return &followRelationService{
		repo: repo,
	}
}

func (f *followRelationService) Follow(ctx context.Context, follower, followee int64) error {
	return f.repo.AddFollowRelation(ctx, domain.FollowRelation{
		Followee: followee,
		Follower: follower,
	})
}

func (f *followRelationService) CancelFollow(ctx context.Context, follower, followee int64) error {
	return f.repo.InactiveFollowRelation(ctx, follower, followee)
}

func (f *followRelationService) GetFollowee(ctx context.Context,
	follower, offset, limit int64) ([]domain.FollowRelation, error) {
	return f.repo.GetFollowee(ctx, follower, offset, limit)
}

func (f *followRelationService) GetFollower(ctx context.Context,
	followee, offset, limit int64) ([]domain.FollowRelation, error) {
	return f.repo.GetFollower(ctx, followee, offset, limit)
}

func (f *followRelationService) FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error) {
	val, err := f.repo.FollowInfo(ctx, follower, followee)
	return val, err
}

func (f *followRelationService) GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	return f.repo.GetFollowStatics(ctx, uid)
}
