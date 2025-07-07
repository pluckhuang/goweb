package job

import (
	"context"
	"sync"
	"time"

	rlock "github.com/gotomicro/redis-lock"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
)

type LockableJob struct {
	L         logger.LoggerV1
	Timeout   time.Duration
	Client    *rlock.Client
	Key       string
	localLock *sync.Mutex
	lock      *rlock.Lock
}

func NewLockableJob(l logger.LoggerV1, client *rlock.Client, key string, timeout time.Duration) *LockableJob {
	return &LockableJob{
		L:         l,
		Client:    client,
		Key:       key,
		Timeout:   timeout,
		localLock: &sync.Mutex{},
	}
}

// 获取分布式锁并自动续约
func (lj *LockableJob) EnsureLock() bool {
	lj.localLock.Lock()
	defer lj.localLock.Unlock()
	if lj.lock != nil {
		return true
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	lock, err := lj.Client.Lock(ctx, lj.Key, lj.Timeout,
		&rlock.FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      3,
		}, time.Second)
	if err != nil {
		lj.L.Warn("获取分布式锁失败", logger.Error(err))
		return false
	}
	lj.lock = lock
	go func() {
		er := lock.AutoRefresh(lj.Timeout/2, lj.Timeout)
		if er != nil {
			lj.localLock.Lock()
			lj.lock = nil
			lj.localLock.Unlock()
		}
	}()
	return true
}

func (lj *LockableJob) Unlock() error {
	lj.localLock.Lock()
	lock := lj.lock
	lj.localLock.Unlock()
	if lock == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return lock.Unlock(ctx)
}
