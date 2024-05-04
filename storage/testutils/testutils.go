package testutils

import (
	"database/sql"
	"testing"

	"github.com/garnizeH/dimdim/storage/repo"

	_ "github.com/mattn/go-sqlite3"
)

func Queries(t *testing.T) (*repo.Queries, func() error) {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to create a sqlite database in memory for tests: %v", err)
	}

	if err := repo.Migrate(db); err != nil {
		t.Fatalf("failed to migrate the sqlite database in memory for tests: %v", err)
	}

	return repo.New(db), db.Close
}
