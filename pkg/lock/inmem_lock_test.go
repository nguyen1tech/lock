package lock

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestAcquire(t *testing.T) {
	type testcase struct {
		callRelease bool
		wantAcquire bool
	}

	var testcases = []testcase{
		{
			callRelease: false,
			wantAcquire: true,
		}, {
			callRelease: false,
			wantAcquire: false,
		}, {
			callRelease: true,
			wantAcquire: true,
		},
	}

	storage := NewInMemoryStorage()
	locker := NewInMemLocker(storage, "instance-01", 5*time.Second)
	for _, tc := range testcases {
		ctx := context.TODO()
		if tc.callRelease {
			_ = locker.Release(ctx)
		}
		err := locker.Acquire(ctx)
		gotAcquired := err == nil
		if gotAcquired != tc.wantAcquire {
			t.Errorf("want acquiried: %+v but got: %+v", tc.wantAcquire, gotAcquired)
		}
	}
}

func TestRelease(t *testing.T) {
	type testcase struct {
		callAcquire bool
		wantAcquire bool
	}

	var testcases = []testcase{
		{
			callAcquire: false,
			wantAcquire: false,
		}, {
			callAcquire: true,
			wantAcquire: true,
		},
	}
	storage := NewInMemoryStorage()
	locker := NewInMemLocker(storage, "instance-01", 5*time.Second)
	for _, tc := range testcases {
		ctx := context.TODO()
		if tc.callAcquire {
			_ = locker.Acquire(ctx)
		}
		err := locker.Release(ctx)
		gotAcquired := err == nil
		if gotAcquired != tc.wantAcquire {
			t.Errorf("want acquire: %+v but got: %+v", tc.wantAcquire, gotAcquired)
		}
	}
}

func TestTTLExpiration(t *testing.T) {
	storage := NewInMemoryStorage()
	locker := NewInMemLocker(storage, "instance-01", 100*time.Millisecond)

	ctx := context.TODO()
	err := locker.Acquire(ctx)
	if err != nil {
		t.Errorf("lock should be acquired but got err: %+v", err)
	}

	err = locker.Acquire(ctx)
	if err == nil {
		t.Errorf("lock should not be acquired b/c the expiration time is not reached")
	}

	time.Sleep(200 * time.Millisecond)
	err = locker.Acquire(ctx)
	if err != nil {
		t.Errorf("lock should be acquired after the expiration time is reached")
	}
}

func TestAcquireAndReleaseConcurrently(t *testing.T) {
	storage := NewInMemoryStorage()
	locker := NewInMemLocker(storage, "instance-01", 5*time.Second)
	num := 100
	resCh := make(chan bool, num)
	var wg sync.WaitGroup
	fn := func() {
		defer wg.Done()
		ctx := context.TODO()
		err := locker.Acquire(ctx)
		if err != nil {
			resCh <- false
			return
		}
		resCh <- true
	}

	for i := 0; i < num; i++ {
		wg.Add(1)
		go fn()
	}

	wg.Wait()

	acquiredCount := 0
	notAcquiredCount := 0
	for i := 0; i < num; i++ {
		value := <-resCh
		if value {
			acquiredCount++
		} else {
			notAcquiredCount++
		}
	}
	if acquiredCount != 1 {
		t.Errorf("want 1 out of: %d goroutines to acquire the lock but got: %d", num, acquiredCount)
	}
	if notAcquiredCount != num-acquiredCount {
		t.Errorf("want %d out of: %d goroutines to not acquire the lock but got: %d", num-acquiredCount, num, notAcquiredCount)
	}
}

func BenchmarkAcquireAndRelease(b *testing.B) {
	storage := NewInMemoryStorage()
	locker := NewInMemLocker(storage, "instance-01", 1*time.Second)
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_ = locker.Acquire(ctx)
		_ = locker.Release(ctx)
	}
}
