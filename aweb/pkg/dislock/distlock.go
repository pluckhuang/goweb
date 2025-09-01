package distlock

import (
	"context"
	"errors"
	"time"

	"github.com/go-redsync/redsync/v4"
	redsync_goredis "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"

	clientv3 "go.etcd.io/etcd/client/v3"
	concurrency "go.etcd.io/etcd/client/v3/concurrency"
)

// 分布式锁统一接口
type DistributedLock interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (UnlockFunc, error)
}

// 解锁函数
type UnlockFunc func() error

// ================= Redis 实现 =================

type RedisLock struct {
	rs *redsync.Redsync
}

func NewRedisLock(client *redis.Client) *RedisLock {
	pool := redsync_goredis.NewPool(client)
	return &RedisLock{
		rs: redsync.New(pool),
	}
}

func (r *RedisLock) Lock(ctx context.Context, key string, ttl time.Duration) (UnlockFunc, error) {
	mutex := r.rs.NewMutex(key, redsync.WithExpiry(ttl))
	if err := mutex.LockContext(ctx); err != nil {
		return nil, err
	}
	return func() error { _, err := mutex.Unlock(); return err }, nil
}

// ================= Etcd 实现 =================

type EtcdLock struct {
	client *clientv3.Client
}

func NewEtcdLock(client *clientv3.Client) *EtcdLock {
	return &EtcdLock{client: client}
}

func (e *EtcdLock) Lock(ctx context.Context, key string, ttl time.Duration) (UnlockFunc, error) {
	s, err := concurrency.NewSession(e.client, concurrency.WithTTL(int(ttl.Seconds())))
	if err != nil {
		return nil, err
	}
	m := concurrency.NewMutex(s, key)
	if err := m.Lock(ctx); err != nil {
		return nil, err
	}
	return func() error {
		defer s.Close()
		return m.Unlock(ctx)
	}, nil
}

// ================= 工厂类 =================

type LockType string

const (
	RedisType LockType = "redis"
	EtcdType  LockType = "etcd"
)

type LockManager struct {
	lock DistributedLock
}

func NewLockManagerRedis(client *redis.Client) *LockManager {
	return &LockManager{lock: NewRedisLock(client)}
}

func NewLockManagerEtcd(client *clientv3.Client) *LockManager {
	return &LockManager{lock: NewEtcdLock(client)}
}

func (lm *LockManager) Acquire(ctx context.Context, key string, ttl time.Duration) (UnlockFunc, error) {
	if lm.lock == nil {
		return nil, errors.New("lock backend not initialized")
	}
	return lm.lock.Lock(ctx, key, ttl)
}
