package shorten

import (
	"github.com/poggerr/go_shortener/internal/app/storage"
	"math/rand"
)

func Shorting(oldUrl string) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	ans := storage.SaveToMap(string(b), oldUrl)
	return ans
}

func UnShoring(newUrl string) string {
	ans := storage.GetFromMap(newUrl)
	return ans
}
