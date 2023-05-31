package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type Storage struct {
	Data map[string]string
	path string
}

func NewStorage(p string) *Storage {
	return &Storage{
		Data: make(map[string]string),
		path: p,
	}
}

func (strg *Storage) Save(key, value string) (string, error) {
	strg.ReadFromFile()
	_, ok := strg.Data[key]
	if ok {
		return "", errors.New("Hey")
	}
	strg.Data[key] = value
	strg.SaveToFile()
	return key, nil
}

func (strg *Storage) OldUrl(key string) (string, error) {
	strg.ReadFromFile()
	val, ok := strg.Data[key]
	if !ok {
		return "", errors.New("Такой ссылки нет. Введите запрос повторно")
	}
	return val, nil
}

func (strg *Storage) SaveToFile() {
	file, err := os.OpenFile(strg.path, os.O_WRONLY|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	data, err := json.Marshal(strg.Data)
	file.Write(data)
}

func (strg *Storage) ReadFromFile() {
	file, err := os.OpenFile(strg.path, os.O_RDONLY|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = json.Unmarshal(data, &strg.Data)
	if err != nil {
		fmt.Println(err.Error())
	}
}
