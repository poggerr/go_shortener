package storage

import (
	"errors"
)

type Storage struct {
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{data: make(map[string]string)} // создаем стурктуру и возвращаем ссылку на нее
}

func (strg *Storage) Save(key, value string) (string, error) {
	_, ok := strg.data[key]
	if ok {
		return "", errors.New("Повторите запрос")
	}
	strg.data[key] = value
	return key, nil
}

func (strg *Storage) OldUrl(key string) (string, error) {
	val, ok := strg.data[key]
	if !ok {
		return "", errors.New("Такой ссылки нет. Введите запрос повторно")
	}
	return val, nil
}
