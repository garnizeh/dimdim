package datastore

import (
	"embed"

	"github.com/garnizeH/dimdim/storage"
)

//go:embed sql/migrations/*
var Migrations embed.FS

func Factory(tx storage.DBTX) *Queries {
	return New(tx)
}
