package storage

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/poggerr/go_shortener/internal/logger"
)

// GetUserID получение id пользователя по username
func (strg *Storage) GetUserID(username string) *uuid.UUID {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id *uuid.UUID
	ans := strg.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE username=$1", username)
	errScan := ans.Scan(&id)
	if errScan != nil {
		logger.Initialize().Info(errScan)
	}
	return id
}
