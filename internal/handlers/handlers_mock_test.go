package handlers

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/poggerr/go_shortener/internal/models"
	"github.com/poggerr/go_shortener/internal/service"
)

const mockedID = "1111"

var ErrNotExistedID = errors.New("mocked fail, use id = 1111 to get stored link")

// RepoMock - простейший мок для Repository интерфейса
type RepoMock struct {
	singleItemStorage string
}

var _ service.URLShortenerService = (*RepoMock)(nil)

func (rm RepoMock) IsExist(_ context.Context, _ string) bool {
	return false
}

func (rm RepoMock) Store(_ context.Context, _ *uuid.UUID, _ string) (id string, err error) {
	// rm.singleItemStorage = link
	return mockedID, nil
}

func (rm RepoMock) Restore(_ context.Context, id string) (link string, err error) {
	if id != mockedID {
		return "", ErrNotExistedID
	}
	return rm.singleItemStorage, nil
}

func (rm RepoMock) Unstore(_ context.Context, _ string, _ []string) {
}

func (rm RepoMock) GetUserStorage(_ context.Context, _ *uuid.UUID, _ string) (map[string]string, error) {
	return map[string]string{mockedID: rm.singleItemStorage}, nil
}

func (rm RepoMock) StoreBatch(_ context.Context, _ *uuid.UUID, _ models.BatchList, _ string) (models.BatchList, error) {
	return models.BatchList{}, nil
}

func (rm RepoMock) Close() error {
	return nil
}

func (rm RepoMock) Ping(_ context.Context) error {
	return nil
}

func (rm RepoMock) Delete(_ context.Context, _ *uuid.UUID, _ []string) {
	return
}
