package service

import (
	"context"
	"fmt"
	"github.com/poggerr/go_shortener/internal/app/models"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/logger"
	"math/rand"
	"time"
)

func Shorting(longURL string, strg *storage.Storage) string {
	if longURL == "" {
		logger.Initialize().Error("Введите ссылку")
		return ""
	}
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	shortURL := string(b)

	url := strg.Save(shortURL, longURL)
	return url
}

func UnShoring(shortURL string, strg *storage.Storage) (string, error) {
	ans, err := strg.LongURL(shortURL)
	return ans, err
}

func SaveMultipleToDB(list models.BatchList, strg *storage.Storage) models.BatchList {
	tx, err := strg.DB.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err != nil {
		logger.Initialize().Error(err)
	}
	for i, v := range list {
		shortUrl := Shorting(v.OriginalUrl, strg)
		list[i].ShortUrl = shortUrl
		query := fmt.Sprintf("INSERT INTO urls (correlation_id, longurl, shorturl) VALUES('%s', '%s', '%s')", v.CorrelationId, v.OriginalUrl, shortUrl)
		_, err = tx.ExecContext(ctx, query)
		if err != nil {
			tx.Rollback()
			logger.Initialize().Info("Ошибка при отправке запроса ", err)
		}
	}
	tx.Commit()
	fmt.Println(list)

	return list

}
