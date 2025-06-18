package idempotent

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
)

// mockHandler 用于模拟业务处理函数
func mockHandler(success bool, dbDup bool) HandlerFunc {
    return func(ctx context.Context) error {
        if dbDup {
            return gorm.ErrDuplicatedKey
        }
        if !success {
            return errors.New("业务处理失败")
        }
        return nil
    }
}

func setupTestRedis() *redis.Client {
    // 使用本地 Redis 测试实例，建议用 Docker 启动临时 redis-bloom 容器
    return redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        DB:   1, // 用测试库，避免污染生产数据
    })
}

func TestIdempotent_Do(t *testing.T) {
    ctx := context.Background()
    rdb := setupTestRedis()
    defer rdb.FlushDB(ctx)

    // 预先创建布隆过滤器
    _, err := rdb.Do(ctx, "BF.RESERVE", "test_bloom", 0.01, 1000).Result()
    if err != nil && err.Error() != "ERR item exists" {
        t.Fatalf("布隆过滤器创建失败: %v", err)
    }

    idem := New(rdb, nil, "test_bloom", 10*time.Second)

    // 1. 第一次请求，业务成功
    err = idem.Do(ctx, "key1", mockHandler(true, false))
    if err != nil {
		t.Errorf("第一次请求应成功, 得到: %v", err)
    }

    // 2. 重复请求，布隆过滤器应拦截
    err = idem.Do(ctx, "key1", mockHandler(true, false))
    if err == nil || err.Error() != "重复请求(布隆过滤器)" {
        t.Errorf("重复请求应被布隆过滤器拦截, 得到: %v", err)
    }

    // 3. 并发请求，SETNX 拦截
    // 先手动设置 Redis key
    ok, _ := rdb.SetNX(ctx, "idempotent:key2", 1, 10*time.Second).Result()
    if !ok {
        t.Fatalf("SETNX 预设失败")
    }
    err = idem.Do(ctx, "key2", mockHandler(true, false))
    if err == nil || err.Error() != "请求处理中或已处理 Redis" {
        t.Errorf("SETNX 并发拦截失败, 得到: %v", err)
    }

    // 4. 业务处理失败
    err = idem.Do(ctx, "key3", mockHandler(false, false))
    if err == nil || err.Error() != "业务处理失败" {
        t.Errorf("业务处理失败未正确返回, 得到: %v", err)
    }

    // 5. 数据库唯一索引冲突
    err = idem.Do(ctx, "key4", mockHandler(true, true))
    if err == nil || err.Error() != "数据库唯一索引冲突，重复请求" {
        t.Errorf("唯一索引冲突未正确返回, 得到: %v", err)
    }
}