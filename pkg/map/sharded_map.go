package shardedmap

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"sync"
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

func (s *ShardedMap) getShard(key any) (*Shard, error) {
	hasher := fnv.New32a()
	_, err := fmt.Fprint(hasher, key)
	if err != nil {
		return nil, fmt.Errorf("failed to hash key: %w", err)
	}
	idx := int(hasher.Sum32()) % shardsNumber
	return s.shards[idx], nil
}

func (s *ShardedMap) Save(_ context.Context, key any, value any) error {
	shard, err := s.getShard(key)
	if err != nil {
		return fmt.Errorf("failed to get shard: %w", err)
	}

	shard.Lock()
	defer shard.Unlock()

	if _, ok := shard.items[key]; ok {
		return ErrAlreadyExists
	}

	shard.items[key] = value

	return nil
}

func (s *ShardedMap) Get(_ context.Context, key any) (any, error) {
	shard, err := s.getShard(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get shard: %w", err)
	}

	shard.RLock()
	defer shard.RUnlock()

	res, ok := shard.items[key]
	if !ok {
		return "", ErrNotFound
	}

	return res, nil
}
