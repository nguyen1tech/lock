package lock

import "context"

type InMemoryStorage struct {
	m map[string]interface{}
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		m: map[string]interface{}{},
	}
}

func newInMemoryStorage(m map[string]interface{}) *InMemoryStorage {
	return &InMemoryStorage{m}
}

func (s *InMemoryStorage) Get(ctx context.Context, key string) (error, interface{}) {
	value, ok := s.m[key]
	if ok {
		return nil, value
	}
	return ErrKeyNotExisted, nil
}

func (s *InMemoryStorage) Set(ctx context.Context, key string, value interface{}) error {
	s.m[key] = value
	return nil
}

func (s *InMemoryStorage) Del(ctx context.Context, key string) error {
	delete(s.m, key)
	return nil
}
