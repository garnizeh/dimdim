package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/garnizeH/dimdim/storage"
)

type SessionManager struct {
	sm *scs.SessionManager
	db *sql.DB
}

func newSessionManager(dsn string) (*SessionManager, error) {
	db, err := storage.NewDBSqlite(dsn)
	if err != nil {
		return nil, err
	}

	const migrate = `
    CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		data BLOB NOT NULL,
		expiry REAL NOT NULL
	);
	CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions(expiry);`

	if _, err = db.Exec(migrate); err != nil {
		return nil, fmt.Errorf("failed to create sessions table: %w", err)
	}

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 24 * time.Hour * 7
	sessionManager.IdleTimeout = 24 * time.Hour
	sessionManager.Cookie.Name = "_s"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	return &SessionManager{
		sm: sessionManager,
		db: db,
	}, nil
}

func (sm *SessionManager) Close() error {
	return sm.db.Close()
}

func (sm *SessionManager) SessionManager() *scs.SessionManager {
	return sm.sm
}