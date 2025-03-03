package storage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/pressly/goose/v3"

	_ "github.com/mattn/go-sqlite3"
)

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type DB[Queries any] struct {
	factory func(tx DBTX) *Queries
	rddb    *sql.DB

	mu   *sync.Mutex
	wrdb *sql.DB
}

const (
	readDSN  = "%s?_journal=wal&_sync=1&_busy_timeout=5000&_cache_size=10000&_txlock=deferred"
	writeDSN = "%s?_journal=wal&_sync=1&_busy_timeout=5000&_cache_size=10000&_txlock=immediate"
)

func NewDBForTest[Queries any](t *testing.T, migrations embed.FS, factory func(tx DBTX) *Queries) *DB[Queries] {
	t.Helper()

	const filename = "test.db"
	db, err := NewDB(filename, migrations, factory)
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close test database: %v", err)
		}

		if err := os.Remove(filename); err != nil {
			t.Fatalf("failed to remove test database %q: %v", filename, err)
		}
	})

	return db
}

func NewDB[Queries any](dsn string, migrations embed.FS, factory func(tx DBTX) *Queries) (*DB[Queries], error) {
	db, err := NewDBSqlite(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open the database connection %q: %w", dsn, err)
	}

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		db.Close()
		return nil, err
	}

	if err := goose.Up(db, "sql/migrations"); err != nil {
		db.Close()
		return nil, err
	}
	db.Close()

	wrdb, err := sql.Open("sqlite3", fmt.Sprintf(writeDSN, dsn))
	if err != nil {
		return nil, err
	}
	wrdb.SetMaxOpenConns(1)

	rddb, err := sql.Open("sqlite3", fmt.Sprintf(readDSN, dsn))
	if err != nil {
		wrdb.Close()
		return nil, err
	}

	return &DB[Queries]{
		factory: factory,
		rddb:    rddb,
		mu:      &sync.Mutex{},
		wrdb:    wrdb,
	}, nil
}

func (db *DB[Queries]) RDBMS() *sql.DB {
	return db.wrdb
}

func (db *DB[Queries]) Close() error {
	return errors.Join(db.rddb.Close(), db.wrdb.Close())
}

func (db *DB[Queries]) Read(ctx context.Context, f func(queries *Queries) error) error {
	return db.transaction(ctx, db.rddb, f)
}

func (db *DB[Queries]) Write(ctx context.Context, f func(queries *Queries) error) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.transaction(ctx, db.wrdb, f)
}

func (db *DB[Queries]) transaction(ctx context.Context, rdbms *sql.DB, f func(queries *Queries) error) error {
	tx, err := rdbms.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin the database transaction: %w", err)
	}

	if err := f(db.factory(tx)); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			err = errors.Join(err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

func NoRows(err error) bool {
	return err != nil && errors.Is(err, sql.ErrNoRows)
}

func NewDBSqlite(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf(writeDSN, dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to open the database connection %q: %w", dsn, err)
	}

	return db, nil
}