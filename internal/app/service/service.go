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

func ServiceCreate(longURL string, strg *storage.Storage, userID string) (string, error) {
	shortURL := Shorting(longURL)
	strg.Save(shortURL, longURL)
	if strg.DB == nil {
		return shortURL, nil
	}
	ans, err := strg.SaveToDB(longURL, shortURL, userID)
	if err != nil {
		return ans, err
	}
	return shortURL, nil
}

func ServiceSaveLocal(longURL string, strg *storage.Storage) string {
	shortURL := Shorting(longURL)
	strg.Save(shortURL, longURL)
	return shortURL
}

func Take(shortURL string, strg *storage.Storage) (string, bool, error) {
	ans, err := strg.LongURL(shortURL)
	if strg.DB != nil {
		isDelete := strg.TakeLongURLIsDelete(shortURL)
		return ans, isDelete, err
	}
	return ans, false, err
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
		shortURL := ServiceSaveLocal(v.OriginalURL, strg)
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

func (r *URLRepo) DeleteAsync(ids []string, userID string) error {
	r.urlsToDeleteChan <- storage.UserURLs{UserID: userID, URLs: ids}
	return nil
}

func (r *URLRepo) WorkerDeleteURLs() {
	for urls := range r.urlsToDeleteChan {
		r.repository.DeleteUrls(urls)
	}
}
