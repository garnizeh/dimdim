package tag_test

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/garnizeH/dimdim/service"
// 	"github.com/garnizeH/dimdim/service/tag"
// 	"github.com/garnizeH/dimdim/storage"
// )

// var (
// 	validTagName = "test"
// )

// func svcTagAndFuncDBClose(t *testing.T) (*tag.Service, func() error) {
// 	t.Helper()

// 	queries, funcDBClose := storage.TestDB()
// 	svcTag := tag.New(queries)

// 	return svcTag, funcDBClose
// }

// func TestService_CreateTag(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		argName string
// 		wantErr error
// 	}{
// 		{
// 			name:    "empty name",
// 			argName: "",
// 			wantErr: service.ErrInvalidParam,
// 		},
// 		{
// 			name:    "valid name",
// 			argName: "valid",
// 			wantErr: nil,
// 		},
// 		{
// 			name:    "repeated name (name must be unique)",
// 			argName: "valid",
// 			wantErr: service.ErrUniqueParam,
// 		},
// 	}

// 	ctx := context.Background()
// 	svcTag, funcDBClose := svcTagAndFuncDBClose(t)
// 	defer func() {
// 		if err := funcDBClose(); err != nil {
// 			t.Errorf("failed to close the sqlite database connection: %v", err)
// 		}
// 	}()

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := svcTag.CreateTag(ctx, tt.argName)
// 			if !errors.Is(err, tt.wantErr) {
// 				t.Errorf("%q got error = %v, want error %v", tt.name, err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestService_DeleteTag(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		argID   int64
// 		wantErr error
// 	}{
// 		{
// 			name:    "zero id",
// 			argID:   0,
// 			wantErr: service.ErrInvalidParam,
// 		},
// 		{
// 			name:    "valid id",
// 			argID:   1,
// 			wantErr: nil,
// 		},
// 		{
// 			name:    "inexistent id (deleted)",
// 			argID:   1,
// 			wantErr: service.ErrNotFound,
// 		},
// 		{
// 			name:    "inexistent id",
// 			argID:   2,
// 			wantErr: service.ErrNotFound,
// 		},
// 	}

// 	ctx := context.Background()
// 	svcTag, funcDBClose := svcTagAndFuncDBClose(t)
// 	defer func() {
// 		if err := funcDBClose(); err != nil {
// 			t.Errorf("failed to close the sqlite database connection: %v", err)
// 		}
// 	}()

// 	if err := svcTag.CreateTag(ctx, validTagName); err != nil {
// 		t.Fatalf("failed to create a tag record %q in the sqlite database: %v", validTagName, err)
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := svcTag.DeleteTag(ctx, tt.argID)
// 			if !errors.Is(err, tt.wantErr) {
// 				t.Errorf("%q got error = %v, want error %v", tt.name, err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestService_GetTagByID(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		argID    int64
// 		wantErr  error
// 		wantName string
// 	}{
// 		{
// 			name:     "zero id",
// 			argID:    0,
// 			wantErr:  service.ErrInvalidParam,
// 			wantName: "",
// 		},
// 		{
// 			name:     "valid id",
// 			argID:    1,
// 			wantErr:  nil,
// 			wantName: validTagName,
// 		},
// 		{
// 			name:     "inexistent id",
// 			argID:    2,
// 			wantErr:  service.ErrNotFound,
// 			wantName: "",
// 		},
// 	}

// 	ctx := context.Background()
// 	svcTag, funcDBClose := svcTagAndFuncDBClose(t)
// 	defer func() {
// 		if err := funcDBClose(); err != nil {
// 			t.Errorf("failed to close the sqlite database connection: %v", err)
// 		}
// 	}()

// 	if err := svcTag.CreateTag(ctx, validTagName); err != nil {
// 		t.Fatalf("failed to create a tag record %q in the sqlite database: %v", validTagName, err)
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := svcTag.GetTagByID(ctx, tt.argID)
// 			if !errors.Is(err, tt.wantErr) {
// 				t.Errorf("%q got error = %v, want error %v", tt.name, err, tt.wantErr)
// 			}
// 			if err != nil {
// 				return
// 			}

// 			if got.Name != tt.wantName {
// 				t.Errorf("%q got name = %q, want name %q", tt.name, got.Name, tt.wantName)
// 			}
// 		})
// 	}
// }

// func TestService_GetTagByName(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		argName  string
// 		wantErr  error
// 		wantName string
// 	}{
// 		{
// 			name:     "empty name",
// 			argName:  "",
// 			wantErr:  service.ErrInvalidParam,
// 			wantName: "",
// 		},
// 		{
// 			name:     "valid name",
// 			argName:  validTagName,
// 			wantErr:  nil,
// 			wantName: validTagName,
// 		},
// 		{
// 			name:     "inexistent name",
// 			argName:  "inexistent",
// 			wantErr:  service.ErrNotFound,
// 			wantName: "",
// 		},
// 	}

// 	ctx := context.Background()
// 	svcTag, funcDBClose := svcTagAndFuncDBClose(t)
// 	defer func() {
// 		if err := funcDBClose(); err != nil {
// 			t.Errorf("failed to close the sqlite database connection: %v", err)
// 		}
// 	}()

// 	if err := svcTag.CreateTag(ctx, validTagName); err != nil {
// 		t.Fatalf("failed to create a tag record %q in the sqlite database: %v", validTagName, err)
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := svcTag.GetTagByName(ctx, tt.argName)
// 			if !errors.Is(err, tt.wantErr) {
// 				t.Errorf("%q got error = %v, want error %v", tt.name, err, tt.wantErr)
// 			}
// 			if err != nil {
// 				return
// 			}

// 			if got.Name != tt.wantName {
// 				t.Errorf("%q got name = %q, want name %q", tt.name, got.Name, tt.wantName)
// 			}
// 		})
// 	}
// }

// func TestService_ListAllTags(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		seedNames []string
// 		wantErr   error
// 		wantNames []string
// 	}{
// 		{
// 			name:      "empty tags",
// 			seedNames: []string{},
// 			wantErr:   nil,
// 			wantNames: []string{},
// 		},
// 		{
// 			name:      "one tag",
// 			seedNames: []string{"one"},
// 			wantErr:   nil,
// 			wantNames: []string{"one"},
// 		},
// 		{
// 			name:      "two tags",
// 			seedNames: []string{"one", "two"},
// 			wantErr:   nil,
// 			wantNames: []string{"one", "two"},
// 		},
// 	}

// 	ctx := context.Background()

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			svcTag, funcDBClose := svcTagAndFuncDBClose(t)
// 			defer func() {
// 				if err := funcDBClose(); err != nil {
// 					t.Errorf("%q failed to close the sqlite database connection: %v", tt.name, err)
// 				}
// 			}()

// 			for _, name := range tt.seedNames {
// 				if err := svcTag.CreateTag(ctx, name); err != nil {
// 					t.Fatalf("%q failed to create a tag record %q in the sqlite database: %v", tt.name, name, err)
// 				}
// 			}

// 			got, err := svcTag.ListAllTags(ctx)
// 			if !errors.Is(err, tt.wantErr) {
// 				t.Errorf("%q got error = %v, want error %v", tt.name, err, tt.wantErr)
// 			}
// 			if err != nil {
// 				return
// 			}

// 			if len(got) != len(tt.wantNames) {
// 				t.Errorf("%q got %d names, want %d names", tt.name, len(got), len(tt.wantNames))
// 			}
// 			for i, tag := range got {
// 				if tag.Name != tt.wantNames[i] {
// 					t.Errorf("%q got %d name %q, want name %q", tt.name, i, tag.Name, tt.wantNames[i])
// 				}
// 			}
// 		})
// 	}
// }

// func TestService_UpdateTag(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		seedNames []string
// 		argID     int64
// 		argName   string
// 		wantErr   error
// 		wantNames []string
// 	}{
// 		{
// 			name:      "invalid id",
// 			seedNames: []string{},
// 			argID:     0,
// 			argName:   "",
// 			wantErr:   service.ErrInvalidParam,
// 			wantNames: []string{},
// 		},
// 		{
// 			name:      "invalid name",
// 			seedNames: []string{},
// 			argID:     1,
// 			argName:   "",
// 			wantErr:   service.ErrInvalidParam,
// 			wantNames: []string{},
// 		},
// 		{
// 			name:      "no tags - id not found",
// 			seedNames: []string{},
// 			argID:     1,
// 			argName:   validTagName,
// 			wantErr:   service.ErrNotFound,
// 			wantNames: []string{},
// 		},
// 		{
// 			name:      "one tag - id not found",
// 			seedNames: []string{"one"},
// 			argID:     2,
// 			argName:   validTagName,
// 			wantErr:   service.ErrNotFound,
// 			wantNames: []string{},
// 		},
// 		{
// 			name:      "one tag - valid id",
// 			seedNames: []string{"one"},
// 			argID:     1,
// 			argName:   validTagName,
// 			wantErr:   nil,
// 			wantNames: []string{validTagName},
// 		},
// 		{
// 			name:      "two tags",
// 			seedNames: []string{"one", "two"},
// 			argID:     1,
// 			argName:   validTagName,
// 			wantErr:   nil,
// 			wantNames: []string{validTagName, "two"},
// 		},
// 	}

// 	ctx := context.Background()

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			svcTag, funcDBClose := svcTagAndFuncDBClose(t)
// 			defer func() {
// 				if err := funcDBClose(); err != nil {
// 					t.Errorf("%q failed to close the sqlite database connection: %v", tt.name, err)
// 				}
// 			}()

// 			for _, name := range tt.seedNames {
// 				if err := svcTag.CreateTag(ctx, name); err != nil {
// 					t.Fatalf("%q failed to create a tag record %q in the sqlite database: %v", tt.name, name, err)
// 				}
// 			}

// 			err := svcTag.UpdateTag(ctx, tt.argID, tt.argName)
// 			if !errors.Is(err, tt.wantErr) {
// 				t.Errorf("%q got error = %v, want error %v", tt.name, err, tt.wantErr)
// 			}
// 			if err != nil {
// 				return
// 			}

// 			got, err := svcTag.ListAllTags(ctx)
// 			if err != nil {
// 				t.Errorf("%q got error = %v, want no error", tt.name, err)
// 			}
// 			if len(got) != len(tt.wantNames) {
// 				t.Errorf("%q got %d names, want %d names", tt.name, len(got), len(tt.wantNames))
// 			}
// 			for i, tag := range got {
// 				if tag.Name != tt.wantNames[i] {
// 					t.Errorf("%q got %d name %q, want name %q", tt.name, i, tag.Name, tt.wantNames[i])
// 				}
// 			}
// 		})
// 	}
// }
