package lm

import (
	"time"

	"github.com/redis/go-redis/v9"

	"lock/pkg/lock"
)

type LockManager struct {
	redisClient *redis.Client
	ttl         time.Duration
}

func NewInMemLockManager(ttl time.Duration) *LockManager {
	return &LockManager{
		ttl: ttl,
	}
}

func NewRedisLockManager(redisClient *redis.Client, ttl time.Duration) *LockManager {
	return &LockManager{
		redisClient: redisClient,
		ttl:         ttl,
	}
}

func (d *LockManager) NewInMemoryLocker(name string) *lock.InMemLocker {
	storage := lock.NewInMemoryStorage()
	return lock.NewInMemLocker(storage, name, d.ttl)
}

func (d *LockManager) NewRedisLocker(name string) *lock.RedisLocker {
	return lock.NewRedisLocker(d.redisClient, name, d.ttl)
}
