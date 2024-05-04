package service

import (
	"errors"
	"strings"
)

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
