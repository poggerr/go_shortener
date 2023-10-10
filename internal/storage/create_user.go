package storage

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/poggerr/go_shortener/internal/logger"
)

// CreateUser создание пользователя
func (strg *Storage) CreateUser(username, pass string, id *uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := strg.DB.ExecContext(ctx, "INSERT INTO users (id, username, pass) VALUES ($1, $2, $3)", id, username, pass)
	if err != nil {
		logger.Initialize().Info("Ошибка при создании юзера ", err)
		return err
	}
	return nil
}
