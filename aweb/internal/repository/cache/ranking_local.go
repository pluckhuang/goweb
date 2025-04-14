package cache

import (
	"context"
	"errors"
	"time"

	"github.com/ecodeclub/ekit/syncx/atomicx"
	"github.com/pluckhuang/goweb/aweb/internal/domain"
)

type RankingLocalCache struct {
	data       *atomicx.Value[cacheData]
	expiration time.Duration
}

type cacheData struct {
	topN []domain.Article
	ddl  time.Time
}

func (r *RankingLocalCache) Set(ctx context.Context, arts []domain.Article) error {
	r.data.Store(cacheData{
		topN: arts,
		ddl:  time.Now().Add(r.expiration),
	})
	return nil
}

func (r *RankingLocalCache) Get(ctx context.Context) ([]domain.Article, error) {
	data := r.data.Load()
	if len(data.topN) == 0 || data.ddl.Before(time.Now()) {
		return nil, errors.New("本地缓存失效了")
	}
	return data.topN, nil
}

func (r *RankingLocalCache) ForceGet(ctx context.Context) ([]domain.Article, error) {
	data := r.data.Load()
	if len(data.topN) == 0 {
		return nil, errors.New("本地缓存失效了")
	}
	return data.topN, nil
}
