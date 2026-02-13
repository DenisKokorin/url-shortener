package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	CodeUniqueViolation = "23505"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(url string, alias string) error {
	stmt, err := s.db.Prepare("INSERT INTO url(alias, url) VALUES($1, $2)")
	if err != nil {
		return fmt.Errorf("failed to prepare sql statement: %w", err)
	}

	_, err = stmt.Exec(stmt, alias, url)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == CodeUniqueViolation {
				return storage.ErrURLAlreadyExists
			}
		}

		return fmt.Errorf("failed to insert url to DB: %w", err)
	}

	return nil
}

func (s *Storage) GetLongURL(alias string) (string, error) {
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		return "", fmt.Errorf("failed to prepare sql statement: %w", err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("failed to get url from DB: %w", err)
	}

	return resURL, nil
}
