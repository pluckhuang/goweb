package repository

import (
	"context"

	"github.com/pluckhuang/goweb/aweb/article/domain"
	"github.com/pluckhuang/goweb/aweb/cronjob/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type CachedRankingRepository struct {
	redisCache *cache.RankingRedisCache
	localCache *cache.RankingLocalCache
}

func NewCachedRankingRepository(redisCache *cache.RankingRedisCache, localCache *cache.RankingLocalCache) RankingRepository {
	return &CachedRankingRepository{redisCache: redisCache, localCache: localCache}
}

func (repo *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	res, err := repo.localCache.Get(ctx)
	if err == nil {
		return res, nil
	}
	res, err = repo.redisCache.Get(ctx)
	if err != nil {
		return repo.localCache.ForceGet(ctx)
	}
	_ = repo.localCache.Set(ctx, res)
	return res, nil
}

func (repo *CachedRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	_ = repo.localCache.Set(ctx, arts)
	return repo.redisCache.Set(ctx, arts)
}
