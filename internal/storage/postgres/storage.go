package postgres

import (
	"database/sql"
	"fmt"
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

func (s *Storage) SaveURL(url string, alias string) error {}

func (s *Storage) GetLongURL(alias string) (string, error) {}
