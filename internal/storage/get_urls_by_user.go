package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/poggerr/go_shortener/internal/logger"
	"github.com/poggerr/go_shortener/internal/models"
)

func (strg *Storage) GetUrlsByUserID(id *uuid.UUID, defURL string) (*models.Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := strg.DB.QueryContext(ctx, "SELECT * FROM urls WHERE user_id=$1", id)
	if err != nil {
		logger.Initialize().Info(err)
		return nil, err
	}

	storage := make(models.Storage, 0)
	for rows.Next() {
		var url models.Urls
		if err = rows.Scan(&url.UserID, &url.LongURL, &url.ShortURL, &url.DeletedFlag); err != nil {
			logger.Initialize().Info(err)
		}
		url.ShortURL = defURL + "/" + url.ShortURL
		storage = append(storage, url)
	}

	if err = rows.Err(); err != nil {
		logger.Initialize().Info(err)
	}

	if len(storage) < 1 {
		return nil, errors.New("у пользователя нет ссылок")
	}

	return &storage, nil
}
