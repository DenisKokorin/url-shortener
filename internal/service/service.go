package urlshortenerservice

import (
	"context"
	"errors"
	"log/slog"
	"url-shortener/internal/storage"
	er "url-shortener/pkg/errors"
	"url-shortener/pkg/logger"
)

const (
	retryCount = 10
)

var (
	ErrURLNotFound = errors.New("url not found")
)

type AliasGenerator interface {
	Generate() (string, error)
}

type URLShortenerService struct {
	storage   storage.Storage
	generator AliasGenerator
	log       *slog.Logger
}

func NewURLShortenerService(log *slog.Logger, storage storage.Storage, generator AliasGenerator) *URLShortenerService {
	return &URLShortenerService{
		storage:   storage,
		generator: generator,
		log:       log,
	}
}

func (s *URLShortenerService) GetShortURL(ctx context.Context, url string) (string, error) {
	var alias string

	s.log.Info("saving url", slog.String("url", url))

	for range retryCount {
		var err error

		s.log.Info("generating alias")

		alias, err = s.generator.Generate()
		if err != nil {

			s.log.Error("failed to generate alias", logger.ErrorLog(err))

			return "", err
		}

		s.log.Info("got alias", slog.String("alias", alias))

		err = s.storage.SaveURL(ctx, url, alias)
		if err != nil {
			if errors.Is(err, er.ErrURLAlreadyExists) {

				s.log.Warn("alias already exists, trying get new")

				continue
			}

			s.log.Error("failed to save url", logger.ErrorLog(err))
		}

		return alias, err
	}

	s.log.Info("url added", slog.String("url", url), slog.String("alias", alias))

	return alias, nil
}

func (s *URLShortenerService) GetLongURL(ctx context.Context, alias string) (string, error) {

	s.log.Info("trying get url from storage")

	url, err := s.storage.GetLongURL(ctx, alias)
	if errors.Is(err, er.ErrURLNotFound) {

		s.log.Warn("url not found", slog.String("alias", alias))

		return "", err
	}

	if err != nil {

		s.log.Error("failed to get long url", logger.ErrorLog(err))

		return "", nil
	}

	s.log.Info("got url from storage", slog.String("url", url), slog.String("alias", alias))

	return url, nil
}
