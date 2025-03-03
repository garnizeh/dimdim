package service

import (
	"errors"
	"strings"

	"github.com/garnizeH/dimdim/pkg/argon2id"
	"github.com/garnizeH/dimdim/pkg/mailer"
	"github.com/garnizeH/dimdim/service/user"
	"github.com/garnizeH/dimdim/storage"
	"github.com/garnizeH/dimdim/storage/datastore"
)

type Service struct {
	user *user.Service
}

func New(
	argon *argon2id.Argon2idHash,
	mailer *mailer.Mailer,
	db *storage.DB[datastore.Queries],
) *Service {
	user := user.New(argon, mailer, db)

	return &Service{
		user: user,
	}
}

func (s *Service) User() *user.Service {
	return s.user
}

var (
	ErrInvalidParam = errors.New("invalid param")
	ErrUniqueParam  = errors.New("param violated unique constraint")
	ErrNotFound     = errors.New("found no record")
)

func CheckErr(err error) error {
	if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
		return errors.Join(err, ErrUniqueParam)
	}

	return err
}
