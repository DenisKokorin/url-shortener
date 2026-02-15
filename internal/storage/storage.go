package storage

import (
	"context"
	"fmt"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/memory"
	"url-shortener/internal/storage/postgres"
	"url-shortener/pkg/utils"
)

const (
	PostgresStorage = "postgres"
	MemoryStorage   = "memory"
)

type Storage interface {
	SaveURL(ctx context.Context, url string, alias string) error
	GetLongURL(ctx context.Context, alias string) (string, error)
}

func GetStorageFromConfig(cfg *config.Config) (Storage, error) {
	switch cfg.Storage {
	case PostgresStorage:
		postgresPath := utils.MustGetPostgresPath()
		storage, err := postgres.NewPostgresStorage(postgresPath)
		if err != nil {
			return nil, fmt.Errorf("failed to init postgres storage: %w", err)
		}
		return storage, nil

	case MemoryStorage:
		return memory.NewMemoryStorage(), nil
	default:
		return nil, fmt.Errorf("unknown storage")
	}
}
