package idempotent

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Idempotent struct {
	RedisClient *redis.Client
	DB          *gorm.DB
	BloomKey    string
	Expire      time.Duration
}

// New 创建幂等控制器
func New(redisClient *redis.Client, db *gorm.DB, bloomKey string, expire time.Duration) *Idempotent {
	return &Idempotent{
		RedisClient: redisClient,
		DB:          db,
		BloomKey:    bloomKey,
		Expire:      expire,
	}
}

// HandlerFunc 业务处理函数签名
type HandlerFunc func(ctx context.Context) error

// Do 幂等处理主流程
func (i *Idempotent) Do(ctx context.Context, idempotentKey string, handler HandlerFunc) error {
	// 1. 布隆过滤器预判
	exists, err := i.RedisClient.Do(ctx, "BF.EXISTS", i.BloomKey, idempotentKey).Bool()
	if err == nil && exists {
		return errors.New("重复请求(布隆过滤器)")
	}

	// 2. Redis SETNX 原子写入
	ok, err := i.RedisClient.SetNX(ctx, "idempotent:"+idempotentKey, 1, i.Expire).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("请求处理中或已处理 Redis")
	}

	// 3. 业务处理（唯一索引兜底）
	err = handler(ctx)
	if err != nil {
		// 数据库唯一索引冲突
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("数据库唯一索引冲突，重复请求")
		}
		return err
	}

	// 4. 处理成功后，写入布隆过滤器
	_, _ = i.RedisClient.Do(ctx, "BF.ADD", i.BloomKey, idempotentKey).Result()
	return nil
}
