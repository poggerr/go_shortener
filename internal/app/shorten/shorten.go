package shorten

import (
	"errors"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"math/rand"
)

func Shorting(oldUrl string, strg storage.Storage) (string, error) {
	if oldUrl == "" {
		err := errors.New("Введите ссылку")
		return "", err
	}
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	ans := strg.Save(string(b), oldUrl)
	return ans, nil
}

func UnShoring(newUrl string, strg storage.Storage) string {
	ans := strg.OldUrl(newUrl)
	return ans
}
