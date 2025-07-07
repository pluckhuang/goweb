package cache

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"time"

	"github.com/pluckhuang/goweb/aweb/interactive/domain"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
)

const (
	fieldReadCnt    = "read_cnt"
	fieldLikeCnt    = "like_cnt"
	fieldCollectCnt = "collect_cnt"
	ErrKeyNotExist  = redis.Nil
)

type InteractiveCache interface {
	// IncrReadCntIfPresent 如果在缓存中有对应的数据，就 +1
	IncrReadCntIfPresent(ctx context.Context,
		biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context,
		biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context,
		biz string, bizId int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
	// Get 查询缓存中数据
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error
}

type InteractiveRedisCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewInteractiveRedisCache(client redis.Cmdable) InteractiveCache {
	return &InteractiveRedisCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (r *InteractiveRedisCache) IncrReadCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)},
		fieldReadCnt, 1).Err()
}

func (r *InteractiveRedisCache) IncrLikeCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)},
		fieldLikeCnt, 1).Err()
}

func (r *InteractiveRedisCache) DecrLikeCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)},
		fieldLikeCnt, -1).Err()
}

func (r *InteractiveRedisCache) IncrCollectCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)},
		fieldCollectCnt, 1).Err()
}

func (r *InteractiveRedisCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}

func (r *InteractiveRedisCache) Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error {
	key := r.key(biz, bizId)
	err := r.client.HMSet(ctx, key,
		fieldLikeCnt, intr.LikeCnt,
		fieldCollectCnt, intr.CollectCnt,
		fieldReadCnt, intr.ReadCnt).Err()
	if err != nil {
		return err
	}
	return r.client.Expire(ctx, key, r.expiration).Err()
}

func (r *InteractiveRedisCache) Get(ctx context.Context,
	biz string, bizId int64) (domain.Interactive, error) {
	// 直接使用 HMGet，即便缓存中没有对应的 key，也不会返回 error
	data, err := r.client.HGetAll(ctx, r.key(biz, bizId)).Result()
	if err != nil {
		return domain.Interactive{}, err
	}

	if len(data) == 0 {
		// 缓存不存在
		return domain.Interactive{}, ErrKeyNotExist
	}

	// 理论上来说，这里不可能有 error
	collectCnt, _ := strconv.ParseInt(data[fieldCollectCnt], 10, 64)
	likeCnt, _ := strconv.ParseInt(data[fieldLikeCnt], 10, 64)
	readCnt, _ := strconv.ParseInt(data[fieldReadCnt], 10, 64)

	return domain.Interactive{
		Biz:        biz,
		BizId:      bizId,
		CollectCnt: collectCnt,
		LikeCnt:    likeCnt,
		ReadCnt:    readCnt,
	}, err
}
