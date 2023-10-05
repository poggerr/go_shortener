package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/poggerr/go_shortener/internal/models"
	"github.com/poggerr/go_shortener/internal/serviceCreateShortURL"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/poggerr/go_shortener/internal/logger"
)

// SaveToDB сохранение ссылки в базу
func (strg *Storage) SaveToDB(longurl, shorturl string, userID *uuid.UUID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := strg.DB.ExecContext(ctx, "INSERT INTO urls (long_url, short_url, user_id) VALUES ($1, $2, $3)", longurl, shorturl, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			ans := strg.DB.QueryRowContext(ctx, "SELECT short_url FROM urls WHERE long_url=$1", longurl)
			errScan := ans.Scan(&shorturl)
			if errScan != nil {
				logger.Initialize().Info(errScan)
			}
			return shorturl, err
		}
	}
	return "", err
}

// SaveMultipleToDB сохранение списка ссылок
func (strg *Storage) SaveMultipleToDB(list models.BatchList, defURL string) models.BatchList {
	tx, err := strg.DB.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err != nil {
		logger.Initialize().Error(err)
	}
	for i, v := range list {
		shortURL := serviceCreateShortURL.CreateShortURL(v.OriginalURL)
		strg.Save(shortURL, v.OriginalURL)
		list[i].ShortURL = defURL + "/" + shortURL
		query := fmt.Sprintf("INSERT INTO urls (long_url, short_url) VALUES('%s', '%s')", v.OriginalURL, shortURL)
		_, err = tx.ExecContext(ctx, query)
		if err != nil {
			tx.Rollback()
			logger.Initialize().Info("Ошибка при отправке запроса ", err)
		}
	}
	tx.Commit()
	return list
}
