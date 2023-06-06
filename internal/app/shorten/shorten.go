package shorten

import (
	"errors"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"math/rand"
)

func Shorting(oldUrl string, strg *storage.Storage) (string, error) {
	if oldUrl == "" {
		err := errors.New("Введите ссылку")
		return "", err
	}
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	shortUrl := string(b)
	_, err := strg.Save(shortUrl, oldUrl)
	if err != nil {
		return "", err
	}
	return shortUrl, err
}

func UnShoring(newUrl string, strg *storage.Storage) (string, error) {
	ans, err := strg.OldUrl(newUrl)
	return ans, err
}
