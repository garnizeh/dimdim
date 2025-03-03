package argon2id_test

import (
	"errors"
	"testing"

	"github.com/garnizeH/dimdim/pkg/argon2id"
)

func TestArgon2IDHashGenerateHash(t *testing.T) {
	a := argon2id.NewWithDefault()

	tests := []struct {
		name     string
		password []byte
		salt     []byte
		wantErr  error
	}{
		{
			name:     "invalid password",
			password: nil,
			salt:     nil,
			wantErr:  argon2id.ErrInvalidPassword,
		},
		{
			name:     "empty password",
			password: []byte{},
			salt:     nil,
			wantErr:  argon2id.ErrInvalidPassword,
		},
		{
			name:     "invalid salt",
			password: []byte("password"),
			salt:     nil,
			wantErr:  nil,
		},
		{
			name:     "empty salt",
			password: []byte("password"),
			salt:     []byte{},
			wantErr:  nil,
		},
		{
			name:     "valid password and salt",
			password: []byte("password"),
			salt:     []byte("salt"),
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := a.GenerateHash(tt.password, nil)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("%q got error = %v, want error = %v", tt.name, err, tt.wantErr)
				return
			}
		})
	}
}

func TestArgon2IDHashCompare(t *testing.T) {
	a := argon2id.NewWithDefault()
	h, err := a.GenerateHash([]byte("password"), nil)
	if err != nil {
		t.Fatalf("failed to generate hash: %v", err)
	}

	tests := []struct {
		name     string
		hash     []byte
		salt     []byte
		password []byte
		wantErr  error
	}{
		{
			name:     "invalid hash",
			hash:     nil,
			salt:     nil,
			password: nil,
			wantErr:  argon2id.ErrInvalidHash,
		},
		{
			name:     "empty hash",
			hash:     []byte{},
			salt:     nil,
			password: nil,
			wantErr:  argon2id.ErrInvalidHash,
		},
		{
			name:     "invalid salt",
			hash:     h.Hash,
			salt:     nil,
			password: nil,
			wantErr:  argon2id.ErrInvalidSalt,
		},
		{
			name:     "empty salt",
			hash:     h.Hash,
			salt:     []byte{},
			password: nil,
			wantErr:  argon2id.ErrInvalidSalt,
		},
		{
			name:     "invalid password",
			hash:     h.Hash,
			salt:     h.Salt,
			password: nil,
			wantErr:  argon2id.ErrInvalidPassword,
		},
		{
			name:     "empty password",
			hash:     h.Hash,
			salt:     h.Salt,
			password: []byte{},
			wantErr:  argon2id.ErrInvalidPassword,
		},
		{
			name:     "wrong password",
			hash:     h.Hash,
			salt:     h.Salt,
			password: []byte("wrong-password"),
			wantErr:  argon2id.ErrPasswordNotMatch,
		},
		{
			name:     "valid password",
			hash:     h.Hash,
			salt:     h.Salt,
			password: []byte("password"),
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := a.Compare(tt.hash, tt.salt, tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("%q got error = %v, want error = %v", tt.name, err, tt.wantErr)
				return
			}
		})
	}
}
