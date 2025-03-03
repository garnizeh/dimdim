-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  email          TEXT NOT NULL PRIMARY KEY,
  name           TEXT NOT NULL,
  password       BLOB NOT NULL,
  salt           BLOB NOT NULL,
  created_at  INTEGER NOT NULL DEFAULT (unixepoch('subsecond') * 1000),
  updated_at  INTEGER NOT NULL DEFAULT (unixepoch('subsecond') * 1000),
  verified_at INTEGER NOT NULL DEFAULT 0,
  deleted_at  INTEGER NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
