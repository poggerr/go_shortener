package storage

import (
	"context"
	"errors"
	"time"

	"github.com/poggerr/go_shortener/internal/logger"
)

// LongURL получение оригинальной ссылки
func (strg *Storage) LongURL(key string) (string, error) {
	val, ok := strg.data[key]
	if !ok {
		return "/", errors.New("такой ссылки нет. Введите запрос повторно")
	}
	return val, nil
}

// CheckLongURLIsDelete проверка на удаление
func (strg *Storage) CheckLongURLIsDelete(shortURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var isDelete bool
	ans := strg.DB.QueryRowContext(ctx, "SELECT is_deleted FROM urls WHERE short_url=$1", shortURL)
	errScan := ans.Scan(&isDelete)
	if errScan != nil {
		logger.Initialize().Info(errScan)
	}
	return isDelete
}
