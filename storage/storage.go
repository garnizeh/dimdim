package storage

import "database/sql"

func SqliteDB(dsn string) (*sql.DB, error) {
	return sql.Open("sqlite3", dsn)
}
