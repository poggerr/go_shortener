package async

import (
	"context"
	"github.com/google/uuid"

	"github.com/poggerr/go_shortener/internal/storage"
)

type URLRepo struct {
	urlsToDeleteChan chan storage.UserURLs
	repository       storage.Storage
}

func NewDeleter(strg *storage.Storage) *URLRepo {
	return &URLRepo{
		urlsToDeleteChan: make(chan storage.UserURLs, 10),
		repository:       *strg,
	}
}

func (r *URLRepo) DeleteAsync(ids []string, userID *uuid.UUID) error {
	r.urlsToDeleteChan <- storage.UserURLs{UserID: userID, URLs: ids}
	return nil
}

// WorkerDeleteURLs воркер удаления ссылок
func (r *URLRepo) WorkerDeleteURLs(ctx context.Context) {
	for urls := range r.urlsToDeleteChan {
		select {
		case <-ctx.Done():
			return
		default:
			r.repository.DeleteUrls(urls)
		}
	}
}
