package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/pluckhuang/goweb/aweb/internal/domain"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
	//go:embed lua/interactive_ranking_incr.lua
	luaRankingIncr string
	//go:embed lua/interactive_ranking_set.lua
	luaRankingSet string
)

var RankingUpdateErr = errors.New("指定的元素不存在")

const fieldReadCnt = "read_cnt"
const fieldLikeCnt = "like_cnt"
const fieldCollectCnt = "collect_cnt"

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, id int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, id int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, id int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, res domain.Interactive) error
	// IncrRankingIfPresent 如果排名数据存在就+1
	IncrRankingIfPresent(ctx context.Context, biz string, bizId int64) error
	// SetRankingScore 如果排名数据不存在就把数据库中读取到的更新到缓存，如果更新过就+1
	SetRankingScore(ctx context.Context, biz string, bizId int64, count int64) error
	LikeTopN(ctx context.Context, biz string, n int64) ([]domain.Interactive, error)
}

type InteractiveRedisCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewInteractiveRedisCache(client redis.Cmdable) InteractiveCache {
	return &InteractiveRedisCache{
		client: client,
	}
}

func (i *InteractiveRedisCache) Set(ctx context.Context,
	biz string, bizId int64,
	res domain.Interactive) error {
	key := i.key(biz, bizId)
	err := i.client.HSet(ctx, key, fieldCollectCnt, res.CollectCnt,
		fieldReadCnt, res.ReadCnt,
		fieldLikeCnt, res.LikeCnt,
	).Err()
	if err != nil {
		return err
	}
	return i.client.Expire(ctx, key, time.Minute*15).Err()
}

func (i *InteractiveRedisCache) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	key := i.key(biz, id)
	res, err := i.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		return domain.Interactive{}, ErrKeyNotExist
	}
	var intr domain.Interactive
	// 这边是可以忽略错误的
	intr.CollectCnt, _ = strconv.ParseInt(res[fieldCollectCnt], 10, 64)
	intr.LikeCnt, _ = strconv.ParseInt(res[fieldLikeCnt], 10, 64)
	intr.ReadCnt, _ = strconv.ParseInt(res[fieldReadCnt], 10, 64)
	return intr, nil
}

func (i *InteractiveRedisCache) IncrCollectCntIfPresent(ctx context.Context,
	biz string, id int64) error {
	key := i.key(biz, id)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldCollectCnt, 1).Err()
}

func (i *InteractiveRedisCache) IncrLikeCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, 1).Err()
}

func (i *InteractiveRedisCache) DecrLikeCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, -1).Err()
}

func (i *InteractiveRedisCache) IncrReadCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	key := i.key(biz, bizId)
	// 不是特别需要处理 res
	//res, err := i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldReadCnt, 1).Int()
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldReadCnt, 1).Err()
}

func (i *InteractiveRedisCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}

func (i *InteractiveRedisCache) topKey(biz string) string {
	return fmt.Sprintf("top:%d:%s", 100, biz)
}

func (i *InteractiveRedisCache) LikeTopN(ctx context.Context, biz string, n int64) ([]domain.Interactive, error) {
	var start int64 = 0
	var end int64 = n - 1
	key := i.topKey(biz)
	res, err := i.client.ZRevRangeWithScores(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}
	interacts := make([]domain.Interactive, 0, n)
	for _, item := range res {
		idStr := item.Member.(string)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue // 忽略错误
		}
		interacts = append(interacts, domain.Interactive{
			BizId:   id,
			LikeCnt: int64(item.Score),
		})
	}
	return interacts, nil
}

func (r *InteractiveRedisCache) IncrRankingIfPresent(ctx context.Context, biz string, bizId int64) error {
	res, err := r.client.Eval(ctx, luaRankingIncr, []string{r.topKey(biz)}, bizId).Result()
	if err != nil {
		return err
	}
	if res.(int64) == 0 {
		return RankingUpdateErr
	}
	return nil
}

func (r *InteractiveRedisCache) SetRankingScore(ctx context.Context, biz string, bizId int64, count int64) error {
	return r.client.Eval(ctx, luaRankingSet, []string{r.topKey(biz)}, bizId, count).Err()
}
