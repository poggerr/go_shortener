package storage

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/poggerr/go_shortener/internal/logger"
)

// UserURLs структура для удаления ссылок
type UserURLs struct {
	UserID *uuid.UUID
	URLs   []string
}

// DeleteUrls удаление ссылок
func (strg *Storage) DeleteUrls(mas UserURLs) {
	tx, err := strg.DB.Begin()
	if err != nil {
		logger.Initialize().Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, m := range mas.URLs {
		_, err = tx.ExecContext(ctx, "UPDATE urls SET is_deleted=true WHERE short_url=$1 AND user_id=$2", m, mas.UserID)
		if err != nil {
			logger.Initialize().Info("Ошибка при удалении", err)
		}
	}
	tx.Commit()
}
