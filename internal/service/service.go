package urlshortenerservice

import (
	"context"
	"errors"
	"log/slog"
	"url-shortener/internal/storage"
	"url-shortener/pkg/logger"
)

const (
	retryCount = 10
)

var (
	ErrURLNotFound = errors.New("url not found")
)

type AliasGenerator interface {
	Generate(s string) string
}

type Storage interface {
	SaveURL(ctx context.Context, url string, alias string) error
	GetLongURL(ctx context.Context, alias string) (string, error)
}

type URLShortenerService struct {
	storage   Storage
	generator AliasGenerator
	log       *slog.Logger
}

func NewURLShortenerService(log *slog.Logger, storage Storage, generator AliasGenerator) *URLShortenerService {
	return &URLShortenerService{
		storage:   storage,
		generator: generator,
		log:       log,
	}
}

func (s *URLShortenerService) GetShortURL(ctx context.Context, url string) (string, error) {
	var alias string

	for range retryCount {
		alias = s.generator.Generate(url)
		err := s.storage.SaveURL(ctx, url, alias)
		if errors.Is(err, storage.ErrURLAlreadyExists) {
			continue
		}

		s.log.Error("failed to save url", logger.ErrorLog(err))

		return "", err
	}

	return alias, nil
}

func (s *URLShortenerService) GetLongURL(ctx context.Context, alias string) (string, error) {
	url, err := s.storage.GetLongURL(ctx, alias)
	if errors.Is(err, storage.ErrURLNotFound) {
		return "", err
	}

	if err != nil {

		s.log.Error("failed to get long url", logger.ErrorLog(err))

		return "", nil
	}

	return url, nil
}
