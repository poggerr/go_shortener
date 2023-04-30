package app

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var mainMap = make(map[string]string)

func Shorting(oldUrl string) string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	mainMap[string(b)] = oldUrl

	return string(b)
}

func UnShorting(newUrl string) string {
	count, ok := mainMap[newUrl]
	if ok {
		return count
	}
	return "Такого ключа нет"

}
