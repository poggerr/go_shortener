package shorten

import (
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/logger"
	"math/rand"
)

func Shorting(oldUrl string, strg *storage.Storage) string {
	if oldUrl == "" {
		logger.Initialize().Error("Введите ссылку")
		return ""
	}
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	shortUrl := string(b)

	url := strg.Save(shortUrl, oldUrl)
	return url
}

func UnShoring(newUrl string, strg *storage.Storage) (string, error) {
	ans, err := strg.OldUrl(newUrl)
	return ans, err
}
