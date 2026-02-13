-- +goose Up
CREATE TABLE IF NOT EXISTS url(
	alias TEXT NOT NULL PRIMARY KEY,
	url TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);

-- +goose Down
DROP TABLE url;