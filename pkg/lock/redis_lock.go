package lock

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotAcquired = errors.New("lock: not acquired")
	ErrLockNotHeld = errors.New("lock: lock not held")
)

type RedisLocker struct {
	redisClient *redis.Client
	ttl         time.Duration
	name        string
	value       string
}

func NewRedisLocker(redisClient *redis.Client, name string, ttl time.Duration) *RedisLocker {
	return &RedisLocker{
		redisClient: redisClient,
		ttl:         ttl,
		name:        name,
		value:       "test",
	}
}

var scriptLock = `if redis.call("EXISTS", KEYS[1]) == 1 then return 0 end return redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])`

// Acquire acquires a lock
//
// Returns ErrNotAcquired error if failed to acquire a lock otherwise, nil is returned
func (rl *RedisLocker) Acquire(ctx context.Context) error {
	status, err := rl.redisClient.Eval(ctx, scriptLock, []string{rl.name}, rl.value, rl.ttl.Milliseconds()).Result()
	if err != nil {
		return err
	}
	if status != "OK" {
		return ErrNotAcquired
	}
	return nil
}

var scriptUnlock = `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) else return 0 end`

// Release releases a lock held by the locker
//
// Return ErrLockNotHeld error if there was no lock held by the locker, otherwise nil is returned
func (rl *RedisLocker) Release(ctx context.Context) error {
	status, err := rl.redisClient.Eval(ctx, scriptUnlock, []string{rl.name}, rl.value).Result()
	if err == redis.Nil {
		return ErrLockNotHeld
	} else if err != nil {
		return err
	}

	if i, ok := status.(int64); !ok || i != 1 {
		return ErrLockNotHeld
	}
	return nil
}
