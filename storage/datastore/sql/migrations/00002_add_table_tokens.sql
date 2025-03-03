-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tokens (
    token           TEXT NOT NULL PRIMARY KEY,
    type            TEXT NOT NULL,
    email           TEXT NOT NULL,
    expires_at   INTEGER NOT NULL,
    deleted_at   INTEGER NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tokens;
-- +goose StatementEnd
