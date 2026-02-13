package shardedmap

import (
	"context"
	"errors"
	"sync"
	"url-shortener/internal/storage"
)

const (
	shardsNumber = 10
)

var (
	ErrAlreadyExists = errors.New("value already exists")
	ErrNotFound      = errors.New("value not found")
)

type ShardedMap struct {
	shards []*Shard
}

type Shard struct {
	sync.RWMutex
	items map[any]any
}

func NewShardedMap() *ShardedMap {
	s := &ShardedMap{
		shards: make([]*Shard, shardsNumber),
	}

	for i := 0; i < shardsNumber; i++ {
		s.shards[i] = &Shard{
			items: make(map[any]any),
		}
	}

	return s
}

func (s *ShardedMap) getShard(key any) *Shard {
	return nil
}

func (s *ShardedMap) Save(_ context.Context, key any, value any) error {
	shard := s.getShard(key)

	shard.Lock()
	defer shard.Unlock()

	if _, ok := shard.items[key]; !ok {
		return ErrAlreadyExists
	}

	shard.items[key] = value

	return nil
}

func (s *ShardedMap) Get(_ context.Context, key any) (any, error) {
	shard := s.getShard(key)

	shard.RLock()
	defer shard.RUnlock()

	url, ok := shard.items[key]
	if !ok {
		return "", storage.ErrURLNotFound
	}

	return url, nil
}
