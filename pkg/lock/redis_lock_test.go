package lock

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

const lockKey = "__TEST_LOCK_KEY__"

var redisOpts = &redis.Options{
	Network: "tcp",
	Addr:    "127.0.0.1:6379",
}

func TestAcquire_Redis(t *testing.T) {
	ctx := context.Background()
	redisClient := redis.NewClient(redisOpts)
	locker := NewRedisLocker(redisClient, lockKey, 100*time.Second)
	defer teardown(t, redisClient)

	// Acquire the lock
	if err := locker.Acquire(ctx); err != nil {
		t.Errorf("want lock acquired but got error: %+v", err)
	}

	// Re-acquire the lock
	if err := locker.Acquire(ctx); err != ErrNotAcquired {
		t.Errorf("want error: %+v but got: %+v", ErrNotAcquired, err)
	}

	if err := locker.Release(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestRelease_Redis(t *testing.T) {
	ctx := context.Background()
	redisClient := redis.NewClient(redisOpts)
	locker := NewRedisLocker(redisClient, lockKey, time.Second)
	defer teardown(t, redisClient)

	if err := locker.Acquire(ctx); err != nil {
		t.Fatal(err)
	}

	if err := locker.Release(ctx); err != nil {
		t.Errorf("want lock released but got error: %+v", err)
	}
}

func Test_Redis_Concurrent(t *testing.T) {
	ctx := context.Background()
	redisClient := redis.NewClient(redisOpts)
	locker := NewRedisLocker(redisClient, lockKey, 5*time.Second)
	defer teardown(t, redisClient)

	numLocks := int32(0)
	numGoroutines := 100
	var wg sync.WaitGroup
	errChan := make(chan error, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			wait := rand.Int63n(int64(10 * time.Millisecond))
			time.Sleep(time.Duration(wait))

			err := locker.Acquire(ctx)
			if err == ErrNotAcquired {
				return
			} else if err != nil {
				errChan <- err
			} else {
				atomic.AddInt32(&numLocks, 1)
			}
		}()
	}
	wg.Wait()

	close(errChan)
	for err := range errChan {
		t.Fatal(err)
	}

	if int(numLocks) != 1 {
		t.Fatalf("want 1 goroutine to acquire a lock but got %v", numLocks)
	}

	if err := locker.Release(ctx); err != nil {
		t.Fatal(err)
	}
}

func teardown(t *testing.T, rc *redis.Client) {
	if err := rc.Del(context.Background(), lockKey).Err(); err != nil {
		t.Fatal(err)
	}
	if err := rc.Close(); err != nil {
		t.Fatal(err)
	}
}
