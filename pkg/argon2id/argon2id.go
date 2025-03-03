package argon2id

import (
	"bytes"
	"crypto/rand"
	"errors"
	"runtime"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidHash      = errors.New("invalid hash")
	ErrInvalidSalt      = errors.New("invalid salt")
	ErrPasswordNotMatch = errors.New("password doesn't match")
)

// HashSalt struct used to store
// generated hash and salt used to
// generate the hash.
type HashSalt struct {
	Hash, Salt []byte
}

type Argon2idHash struct {
	// time represents the number of
	// passed over the specified memory.
	time uint32
	// cpu memory to be used.
	memory uint32
	// threads for parallelism aspect
	// of the algorithm.
	threads uint8
	// keyLen of the generate hash key.
	keyLen uint32
	// saltLen the length of the salt used.
	saltLen uint32
}

// New returns an Argon2idHash.
func New(time, saltLen uint32, memory uint32, threads uint8, keyLen uint32) *Argon2idHash {
	return &Argon2idHash{
		time:    time,
		saltLen: saltLen,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
	}
}

// NewWithDefault returns an Argon2idHash with default config.
func NewWithDefault() *Argon2idHash {
	// We want to use at most half the cpus available and no more than 4.
	threads := min(4, max(1, uint8(runtime.NumCPU()/2)))

	return New(4, 32, 64*1024, threads, 256)
}

// GenerateHash using the password and provided salt.
// If not salt value provided fallback to random value
// generated of a given length.
func (a *Argon2idHash) GenerateHash(password, salt []byte) (*HashSalt, error) {
	if len(password) == 0 {
		return nil, ErrInvalidPassword
	}

	var err error
	// If salt is not provided generate a salt of
	// the configured salt length.
	if len(salt) == 0 {
		salt, err = randomSecret(a.saltLen)
	}
	if err != nil {
		return nil, err
	}
	// Generate hash
	hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keyLen)
	// Return the generated hash and salt used for storage.
	return &HashSalt{Hash: hash, Salt: salt}, nil
}

// Compare generated hash with store hash.
func (a *Argon2idHash) Compare(hash, salt, password []byte) error {
	if len(hash) == 0 {
		return ErrInvalidHash
	}
	if len(salt) == 0 {
		return ErrInvalidSalt
	}
	if len(password) == 0 {
		return ErrInvalidPassword
	}

	// Generate hash for comparison.
	hashSalt, err := a.GenerateHash(password, salt)
	if err != nil {
		return err
	}
	// Compare the generated hash with the stored hash.
	// If they don't match return error.
	if !bytes.Equal(hash, hashSalt.Hash) {
		return ErrPasswordNotMatch
	}
	
	return nil
}

func randomSecret(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
