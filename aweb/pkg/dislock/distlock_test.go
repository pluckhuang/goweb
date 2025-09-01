package distlock

import (
	"context"
	"sync"
	"testing"
	"time"

	// "example.com/hello/distlock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// go test -v ./aweb/pkg/dislock/... -count=1
func TestRedisLockManager(t *testing.T) {
	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	lockManager := NewLockManagerRedis(rdb)

	ctx := context.Background()
	lockKey := "order:12345"

	// 获取锁
	unlock, err := lockManager.Acquire(ctx, lockKey, 5*time.Second)
	require.NoError(t, err, "获取锁失败")
	require.NotNil(t, unlock, "返回的 unlock 函数不应为 nil")

	t.Logf("锁定成功: %s", lockKey)

	// 模拟业务逻辑
	time.Sleep(2 * time.Second)

	// 解锁
	err = unlock()
	require.NoError(t, err, "解锁失败")

	t.Logf("解锁成功: %s", lockKey)

	// 再次加锁，确认释放成功
	unlock2, err := lockManager.Acquire(ctx, lockKey, 5*time.Second)
	require.NoError(t, err, "再次获取锁失败")
	require.NotNil(t, unlock2, "再次获取锁返回的 unlock 函数不应为 nil")

	// 清理
	_ = unlock2()
}

func TestRedisLockManagerConcurrent(t *testing.T) {
	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	lockManager := NewLockManagerRedis(rdb)
	ctx := context.Background()
	lockKey := "order:concurrent"

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// 用于检测临界区是否被同时进入
	var counter int
	var mu sync.Mutex

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()

			unlock, err := lockManager.Acquire(ctx, lockKey, 5*time.Second)
			require.NoError(t, err, "goroutine %d 获取锁失败", idx)
			require.NotNil(t, unlock, "goroutine %d unlock 函数为空", idx)

			// 临界区开始
			mu.Lock()
			if counter != 0 {
				t.Errorf("goroutine %d 进入临界区时 counter=%d，说明锁失效", idx, counter)
			}
			counter++
			time.Sleep(50 * time.Millisecond) // 模拟业务操作
			counter--
			mu.Unlock()
			// 临界区结束

			err = unlock()
			require.NoError(t, err, "goroutine %d 解锁失败", idx)
		}(i)
	}

	wg.Wait()
	t.Log("所有 goroutine 测试完成，锁互斥验证通过")
}
