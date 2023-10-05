package servicecreateshorturl

import (
	"math/rand"

	"github.com/poggerr/go_shortener/internal/logger"
)

// CreateShortURL создание короткой ссылки
func CreateShortURL(longURL string) string {
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
