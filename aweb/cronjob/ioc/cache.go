package ioc

import (
	"time"

	"github.com/pluckhuang/goweb/aweb/cronjob/repository/cache"
	"github.com/redis/go-redis/v9"
)

// 提供 RankingRedisCache
func InitRankingRedisCache(client redis.Cmdable) *cache.RankingRedisCache {
	return cache.NewRankingRedisCache(client)
}

// 提供 RankingLocalCache
func InitRankingLocalCache() *cache.RankingLocalCache {
	return cache.NewRankingLocalCache(5 * time.Minute) // 5分钟过期时间
}
