package app

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var MainMap = make(map[string]string)

func Shorting(oldUrl string) string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	MainMap[string(b)] = oldUrl

	return string(b)
}

func UnShorting(newUrl string) string {
	count, ok := MainMap[newUrl]
	if ok {
		return count
	}
	return "Такого ключа нет"

}
