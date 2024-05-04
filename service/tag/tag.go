package tag

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/garnizeH/dimdim/service"
	"github.com/garnizeH/dimdim/storage/repo"
)

type Service struct {
	queries *repo.Queries
}

func New(queries *repo.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

func (s *Service) CreateTag(ctx context.Context, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return service.ErrInvalidParam
	}

	if err := s.queries.CreateTag(ctx, repo.CreateTagParams{
		Name:      name,
		CreatedAt: timestamp(),
	}); err != nil {
		return service.CheckErr(err)
	}

	return nil
}

func (s *Service) DeleteTag(ctx context.Context, id int64) error {
	if id == 0 {
		return service.ErrInvalidParam
	}

	if _, err := s.queries.GetTagByID(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(err, service.ErrNotFound)
		}

		return err
	}

	return s.queries.DeleteTag(ctx, repo.DeleteTagParams{
		ID:        id,
		DeletedAt: timestamp(),
	})
}

func (s *Service) GetTagByID(ctx context.Context, id int64) (repo.Tag, error) {
	res := repo.Tag{}
	if id == 0 {
		return res, service.ErrInvalidParam
	}

	tag, err := s.queries.GetTagByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, errors.Join(err, service.ErrNotFound)
		}

		return res, err
	}

	return tag, nil
}

func (s *Service) GetTagByName(ctx context.Context, name string) (repo.Tag, error) {
	res := repo.Tag{}
	name = strings.TrimSpace(name)
	if name == "" {
		return res, service.ErrInvalidParam
	}

	tag, err := s.queries.GetTagByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, errors.Join(err, service.ErrNotFound)
		}

		return res, err
	}

	return tag, nil
}

func (s *Service) ListAllTags(ctx context.Context) ([]repo.Tag, error) {
	return s.queries.ListAllTags(ctx)
}

func (s *Service) UpdateTag(ctx context.Context, id int64, name string) error {
	if id == 0 {
		return service.ErrInvalidParam
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return service.ErrInvalidParam
	}

	if _, err := s.queries.GetTagByID(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(err, service.ErrNotFound)
		}

		return err
	}

	return s.queries.UpdateTag(ctx, repo.UpdateTagParams{
		ID:        id,
		Name:      name,
		UpdatedAt: timestamp(),
	})
}

func timestamp() int64 {
	return time.Now().UTC().UnixMicro()
}
