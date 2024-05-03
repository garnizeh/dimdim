package main

import (
	"log/slog"
	"os"

	"github.com/garnizeH/dimdim/storage"
)

func main() {
	if err := run(); err != nil {
		panic("err")
	}
}

func run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("app starting")

	// TODO: get this value from config
	dsn := ":memory:"
	db, err := storage.SqliteDB(dsn)
	if err != nil {
		logger.Error(
			"failed to open the database connection",
			"dsn", dsn,
			"error", err,
		)

		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(
				"failed to close the database connection",
				"error", err,
			)
		}
	}()

	logger.Info("app finished")
	return nil
}