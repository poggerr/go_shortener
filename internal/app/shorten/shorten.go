package shorten

import (
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/logger"
	"math/rand"
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
