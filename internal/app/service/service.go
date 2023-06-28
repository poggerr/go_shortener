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

func ServiceCreate(longURL, defURL string, strg *storage.Storage, userId string) (string, error) {
	shortURL := Shorting(longURL)
	strg.Save(shortURL, longURL)
	shortURL = defURL + "/" + shortURL
	if strg.DB == nil {
		return shortURL, nil
	}
	ans, err := strg.SaveToDB(longURL, shortURL, userId)
	if err != nil {
		return ans, err
	}
	return shortURL, nil
}

func ServiceSaveLocal(longURL, defURL string, strg *storage.Storage) string {
	shortURL := Shorting(longURL)
	strg.Save(shortURL, longURL)
	shortURL = defURL + "/" + shortURL
	return shortURL
}

func Take(shortURL string, strg *storage.Storage) (string, error) {
	ans, err := strg.LongURL(shortURL)
	return ans, err
}

func Shorting(longURL string) string {
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
	return shortURL
}

func SaveMultipleToDB(list models.BatchList, strg *storage.Storage, defURL string) models.BatchList {
	tx, err := strg.DB.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err != nil {
		logger.Initialize().Error(err)
	}
	for i, v := range list {
		shortURL := ServiceSaveLocal(v.OriginalURL, defURL, strg)
		list[i].ShortURL = shortURL
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

func ServiceDelete(keys []string, str string) {
	for i, v := range keys {
		keys[i] = str + "/" + v
	}

}
