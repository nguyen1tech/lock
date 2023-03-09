package lock

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyNotExisted = errors.New("storage: key not existed")
)

type storage interface {
	Get(ctx context.Context, key string) (error, interface{})
	Set(ctx context.Context, key string, value interface{}) error
	Del(ctx context.Context, key string) error
}

type lockValue struct {
	creationTime time.Time
}

type InMemLocker struct {
	mu      sync.Mutex
	name    string
	ttl     time.Duration
	storage storage
}

func NewInMemLocker(storage storage, name string, ttl time.Duration) *InMemLocker {
	return &InMemLocker{
		ttl:     ttl,
		name:    name,
		storage: storage,
	}
}

func (l *InMemLocker) Acquire(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	err, value := l.storage.Get(ctx, l.name)
	if err == ErrKeyNotExisted {
		lvalue := lockValue{
			creationTime: time.Now(),
		}
		return l.storage.Set(ctx, l.name, lvalue)
	}
	// Check for the expiration time
	lvalue, ok := value.(lockValue)
	if !ok {
		return errors.New("invalid lock value")
	}
	if lvalue.creationTime.Add(l.ttl).After(time.Now()) {
		return ErrNotAcquired
	}
	lvalue = lockValue{
		creationTime: time.Now(),
	}
	return l.storage.Set(ctx, l.name, lvalue)
}

func (l *InMemLocker) Release(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	err, _ := l.storage.Get(ctx, l.name)
	if err != nil {
		if err == ErrKeyNotExisted {
			return ErrLockNotHeld
		}
		return err
	}

	return l.storage.Del(ctx, l.name)
}
