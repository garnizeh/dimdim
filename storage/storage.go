package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func SqliteDB(dsn string) (*sql.DB, error) {
	return sql.Open("sqlite3", dsn)
}
