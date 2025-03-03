package tag

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"strings"
// 	"time"

// 	"github.com/garnizeH/dimdim/service"
// 	"github.com/garnizeH/dimdim/storage/datastore"
// )

// type Service struct {
// 	queries *datastore.Queries
// }

// func New(queries *datastore.Queries) *Service {
// 	return &Service{
// 		queries: queries,
// 	}
// }

// func (s *Service) CreateTag(ctx context.Context, name string) error {
// 	name = strings.TrimSpace(name)
// 	if name == "" {
// 		return service.ErrInvalidParam
// 	}

// 	if err := s.queries.CreateTag(ctx, datastore.CreateTagParams{
// 		Name:      name,
// 		CreatedAt: timestamp(),
// 	}); err != nil {
// 		return service.CheckErr(err)
// 	}

// 	return nil
// }

// func (s *Service) DeleteTag(ctx context.Context, id int64) error {
// 	if id == 0 {
// 		return service.ErrInvalidParam
// 	}

// 	if _, err := s.queries.GetTagByID(ctx, id); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return errors.Join(err, service.ErrNotFound)
// 		}

// 		return err
// 	}

// 	return s.queries.DeleteTag(ctx, datastore.DeleteTagParams{
// 		ID:        id,
// 		DeletedAt: timestamp(),
// 	})
// }

// func (s *Service) GetTagByID(ctx context.Context, id int64) (datastore.Tag, error) {
// 	res := datastore.Tag{}
// 	if id == 0 {
// 		return res, service.ErrInvalidParam
// 	}

// 	tag, err := s.queries.GetTagByID(ctx, id)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return res, errors.Join(err, service.ErrNotFound)
// 		}

// 		return res, err
// 	}

// 	return tag, nil
// }

// func (s *Service) GetTagByName(ctx context.Context, name string) (datastore.Tag, error) {
// 	res := datastore.Tag{}
// 	name = strings.TrimSpace(name)
// 	if name == "" {
// 		return res, service.ErrInvalidParam
// 	}

// 	tag, err := s.queries.GetTagByName(ctx, name)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return res, errors.Join(err, service.ErrNotFound)
// 		}

// 		return res, err
// 	}

// 	return tag, nil
// }

// func (s *Service) ListAllTags(ctx context.Context) ([]datastore.Tag, error) {
// 	return s.queries.ListAllTags(ctx)
// }

// func (s *Service) UpdateTag(ctx context.Context, id int64, name string) error {
// 	if id == 0 {
// 		return service.ErrInvalidParam
// 	}
// 	name = strings.TrimSpace(name)
// 	if name == "" {
// 		return service.ErrInvalidParam
// 	}

// 	if _, err := s.queries.GetTagByID(ctx, id); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return errors.Join(err, service.ErrNotFound)
// 		}

// 		return err
// 	}

// 	return s.queries.UpdateTag(ctx, datastore.UpdateTagParams{
// 		ID:        id,
// 		Name:      name,
// 		UpdatedAt: timestamp(),
// 	})
// }

// func timestamp() int64 {
// 	return time.Now().UTC().UnixMicro()
// }
