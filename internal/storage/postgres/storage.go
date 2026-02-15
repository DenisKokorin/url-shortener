package postgres

import (
	"context"
	"errors"
	"fmt"
	er "url-shortener/pkg/errors"
	"url-shortener/pkg/utils"

	_ "github.com/lib/pq"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	CodeUniqueViolation = "23505"
)

type Storage struct {
	db *pgx.Conn
}

func NewPostgresStorage(path string) (*Storage, error) {
	db, err := pgx.Connect(context.Background(), path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	err = utils.ApplyMigration(path)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(ctx context.Context, url string, alias string) error {
	stmt := "INSERT INTO url(alias, url) VALUES($1, $2)"

	_, err := s.db.Exec(ctx, stmt, alias, url)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == CodeUniqueViolation {
				return er.ErrURLAlreadyExists
			}
		}

		return fmt.Errorf("failed to insert url to DB: %w", err)
	}

	return nil
}

func (s *Storage) GetLongURL(ctx context.Context, alias string) (string, error) {
	stmt := "SELECT url FROM url WHERE alias = $1"

	var resURL string

	err := s.db.QueryRow(ctx, stmt, alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", er.ErrURLNotFound
		}

		return "", fmt.Errorf("failed to get url from DB: %w", err)
	}

	return resURL, nil
}
