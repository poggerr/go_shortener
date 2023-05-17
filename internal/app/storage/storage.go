package storage

import (
	"errors"
)

type Storage map[string]string

func NewStorage() Storage {
	arr := make(Storage)
	return arr
}

func (strg *Storage) Save(key, value string) (string, error) {
	(*strg)[key] = value
	return key, nil
}

func (strg Storage) OldUrl(key string) (string, error) {
	val, ok := strg[key]
	if !ok {
		return "", errors.New("Такой ссылки нет. Введите запрос повторно")
	}
	return val, nil
}
