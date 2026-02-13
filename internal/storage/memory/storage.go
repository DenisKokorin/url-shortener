package memory

import (
	"context"
	"errors"
	"url-shortener/internal/storage"
	shardedmap "url-shortener/pkg/map"
)

type MemoryStorage struct {
	storage *shardedmap.ShardedMap
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		storage: shardedmap.NewShardedMap(),
	}
}

func (ms *MemoryStorage) SaveURL(ctx context.Context, url string, alias string) error {
	err := ms.storage.Save(ctx, url, alias)
	if errors.Is(err, shardedmap.ErrAlreadyExists) {
		return storage.ErrURLAlreadyExists
	}

	return nil
}

func (ms *MemoryStorage) GetLongURL(ctx context.Context, alias string) (string, error) {
	url, err := ms.storage.Get(ctx, alias)
	if errors.Is(err, shardedmap.ErrNotFound) {
		return "", storage.ErrURLNotFound
	}

	return url.(string), nil
}
