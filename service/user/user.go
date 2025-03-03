package user

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/garnizeH/dimdim/pkg/argon2id"
	"github.com/garnizeH/dimdim/pkg/mailer"
	"github.com/garnizeH/dimdim/storage"
	"github.com/garnizeH/dimdim/storage/datastore"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

const (
	tokenSignup        = "SIGNUP"
	tokenResetPassword = "RESET_PASSWORD"

	tokenDurationSignup        = time.Hour * 12
	tokenDurationResetPassword = time.Hour * 1
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailInUse          = errors.New("email already in use")
	ErrUserNotVerified     = errors.New("verify your mail box")
	ErrEmailNotFound       = errors.New("email not found")
	ErrUserAlreadyVerified = errors.New("user already verified")
	ErrInvalidToken        = errors.New("invalid token")
)

type Service struct {
	argon     *argon2id.Argon2idHash
	mailer    *mailer.Mailer
	db        *storage.DB[datastore.Queries]
	userCache *sync.Map
}

func New(
	argon *argon2id.Argon2idHash,
	mailer *mailer.Mailer,
	db *storage.DB[datastore.Queries],
) *Service {
	return &Service{
		argon:     argon,
		mailer:    mailer,
		db:        db,
		userCache: &sync.Map{},
	}
}

type User struct {
	ID    string
	Name  string
	Email string
}

func (s *Service) GetUser(ctx context.Context, id string) (User, error) {
	if v, ok := s.userCache.Load(id); ok {
		return v.(User), nil
	}

	var u datastore.User
	if err := s.db.Read(ctx, func(queries *datastore.Queries) error {
		var err error
		u, err = queries.GetUserByID(ctx, id)
		return err
	}); err != nil {
		return User{}, err
	}

	return s.updateCache(u), nil
}

func (s *Service) Signin(ctx context.Context, email, password string) (User, error) {
	var u datastore.User
	err := s.db.Read(ctx, func(queries *datastore.Queries) error {
		var err error
		u, err = queries.GetUserByEmail(ctx, email)
		if err != nil {
			if storage.NoRows(err) {
				return ErrInvalidCredentials
			}

			return err
		}

		if u.VerifiedAt == 0 {
			return ErrUserNotVerified
		}

		return nil
	})
	if err != nil {
		return User{}, err
	}

	if err := s.argon.Compare(u.Password, u.Salt, []byte(password)); err != nil {
		return User{}, ErrInvalidCredentials
	}

	return s.updateCache(u), nil
}

func (s *Service) Signup(
	ctx context.Context,
	baseURL string,
	email string,
	name string,
	password string,
) error {
	timestamp := time.Now().UTC().UnixMilli()

	uid, err := ulid.New(uint64(timestamp), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to create the user id: %w", err)
	}

	id := uid.String()
	if err := s.db.Write(ctx, func(queries *datastore.Queries) error {
		_, err := queries.GetUserByEmail(ctx, email)
		if err == nil {
			return ErrEmailInUse
		}
		if !storage.NoRows(err) {
			return fmt.Errorf("failed to check for the email existence in the database: %w", err)
		}

		hashSalt, err := s.argon.GenerateHash([]byte(password), nil)
		if err != nil {
			return fmt.Errorf("failed to hash the password: %w", err)
		}

		if err := queries.CreateUser(ctx, datastore.CreateUserParams{
			ID:       id,
			Email:    email,
			Name:     name,
			Password: hashSalt.Hash,
			Salt:     hashSalt.Salt,
		}); err != nil {
			return fmt.Errorf("failed to create the user in the database: %w", err)
		}

		token := uuid.New().String()
		mail := mailer.NewMailSignup(baseURL, email, name, token)

		if err := queries.DeleteSignupTokensByEmail(ctx, email); err != nil {
			return fmt.Errorf("failed to delete existing signup tokens for the email %q in the database: %w", email, err)
		}

		expiresAt := time.Now().Add(tokenDurationSignup).UTC().UnixMilli()
		if err := queries.CreateToken(ctx, datastore.CreateTokenParams{
			Token:     token,
			Type:      tokenSignup,
			Email:     email,
			ExpiresAt: expiresAt,
		}); err != nil {
			return fmt.Errorf("failed to create the signup token in the database: %w", err)
		}

		if err := s.mailer.SendMailSignup(mail); err != nil {
			return fmt.Errorf("failed to send the signup confirmation email: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) ResendSignupToken(
	ctx context.Context,
	baseURL string,
	email string,
) error {
	if err := s.db.Write(ctx, func(queries *datastore.Queries) error {
		u, err := queries.GetUserByEmail(ctx, email)
		if err != nil {
			if storage.NoRows(err) {
				return ErrUserNotVerified
			}

			return fmt.Errorf("failed to check for the email existence in the database: %w", err)
		}
		if u.VerifiedAt > 0 {
			return ErrUserAlreadyVerified
		}

		token := uuid.New().String()
		mail := mailer.NewMailSignup(baseURL, email, u.Name, token)

		if err := queries.DeleteSignupTokensByEmail(ctx, email); err != nil {
			return fmt.Errorf("failed to delete existing signup tokens for the email %q in the database: %w", email, err)
		}

		expiresAt := time.Now().Add(tokenDurationSignup).UTC().UnixMilli()
		if err := queries.CreateToken(ctx, datastore.CreateTokenParams{
			Token:     token,
			Type:      tokenSignup,
			Email:     email,
			ExpiresAt: expiresAt,
		}); err != nil {
			return fmt.Errorf("failed to create the signup token in the database: %w", err)
		}

		if err := s.mailer.SendMailSignup(mail); err != nil {
			return fmt.Errorf("failed to send the signup confirmation email: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (datastore.User, error) {
	var user datastore.User
	now := time.Now().UTC().UnixMilli()
	if err := s.db.Read(ctx, func(queries *datastore.Queries) error {
		registeredToken, err := queries.GetSignupTokenNotExpired(ctx, datastore.GetSignupTokenNotExpiredParams{
			Token:     token,
			ExpiresAt: now,
		})
		if err != nil {
			if storage.NoRows(err) {
				return ErrInvalidToken
			}

			return err
		}

		if err := queries.DeleteToken(ctx, registeredToken.Email); err != nil {
			return err
		}

		user, err = queries.SetUserIsVerified(ctx, registeredToken.Email)
		if err != nil {
			return err
		}

		_ = s.updateCache(user)

		return nil
	}); err != nil {
		return datastore.User{}, err
	}

	return user, nil
}

func (s *Service) updateCache(u datastore.User) User {
	user := User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
	s.userCache.LoadOrStore(u.ID, user)

	return user
}
