-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id             TEXT NOT NULL PRIMARY KEY,
  email          TEXT NOT NULL,
  name           TEXT NOT NULL,
  password       BLOB NOT NULL,
  salt           BLOB NOT NULL,
  created_at  INTEGER NOT NULL DEFAULT (unixepoch('subsecond') * 1000),
  updated_at  INTEGER NOT NULL DEFAULT (unixepoch('subsecond') * 1000),
  verified_at INTEGER NOT NULL DEFAULT 0,
  deleted_at  INTEGER NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS users_email_idx;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
