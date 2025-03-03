-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tags (
  id         INTEGER PRIMARY KEY,
  name       TEXT    NOT NULL UNIQUE,
  created_at INTEGER NOT NULL DEFAULT (unixepoch('subsecond') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (unixepoch('subsecond') * 1000),
  deleted_at INTEGER NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_name ON tags (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tags_name;
DROP TABLE IF EXISTS tags;
-- +goose StatementEnd
