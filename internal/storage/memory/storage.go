package memory

import (
	"context"
	"sync"
	"url-shortener/internal/storage"
)

const (
	shardsNumber = 10
)

type Storage struct {
	shards []*Shard
}

type Shard struct {
	sync.RWMutex
	items map[string]string
}

func NewMemoryStorage() *Storage {
	s := &Storage{
		shards: make([]*Shard, shardsNumber),
	}

	for i := 0; i < shardsNumber; i++ {
		s.shards[i] = &Shard{
			items: make(map[string]string),
		}
	}

	return s
}

func (s *Storage) getShard(key string) *Shard {
	return nil
}

func (s *Storage) SaveURL(_ context.Context, url string, alias string) error {
	shard := s.getShard(alias)

	shard.Lock()
	defer shard.Unlock()

	if _, ok := shard.items[alias]; !ok {
		return storage.ErrURLAlreadyExists
	}

	shard.items[alias] = url

	return nil
}

func (s *Storage) GetLongURL(_ context.Context, alias string) (string, error) {
	shard := s.getShard(alias)

	shard.RLock()
	defer shard.RUnlock()

	url, ok := shard.items[alias]
	if !ok {
		return "", storage.ErrURLNotFound
	}

	return url, nil
}
