package storage

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/garnizeH/dimdim/storage/repo"
	
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrOpenConnection = errors.New("failed to open database connection")
	ErrMigrate        = errors.New("failed to migrate database")
)

func SqliteDB(dsn string) (*sql.DB, error) {
	slog.Info(
		"opening database connection",
		"dsn", dsn,
	)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, errors.Join(err, ErrOpenConnection)
	}

	slog.Info("migrating database")
	if err := repo.Migrate(db); err != nil {
		return nil, errors.Join(err, ErrMigrate)
	}

	return db, nil
}
