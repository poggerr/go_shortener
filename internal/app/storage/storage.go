package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/poggerr/go_shortener/internal/logger"
	"os"
	"path"
)

type LongUrl string

type Url struct {
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

type Storage struct {
	data map[string]string
	path string
}

func NewStorage(p string) *Storage {
	return &Storage{
		data: make(map[string]string),
		path: p,
	}
}

func (strg *Storage) Save(key, value string) string {
	strg.data[key] = value
	if strg.path != "" {
		strg.SaveToFile()
	}
	return key
}

func (strg *Storage) OldUrl(key string) (string, error) {
	val, ok := strg.data[key]
	if !ok {
		return "/", errors.New("Такой ссылки нет. Введите запрос повторно")
	}
	return val, nil
}

func (strg *Storage) SaveToFile() {
	file, err := os.OpenFile(strg.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	data, err := json.Marshal(strg.data)

	data = append(data, '\n')

	_, err = file.Write(data)
	if err != nil {
		logger.Log.Error("Ошибка при сохранении файла ", err)
	}
}

func (strg *Storage) RestoreFromFile() {
	dir, _ := path.Split(strg.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0666)
		if err != nil {
			logger.Initialize().Info(err)
		}
	}
	file, err := os.OpenFile(strg.path, os.O_RDONLY|os.O_CREATE, 0666)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Initialize().Error(err)
		}
	}(file)
	if err != nil {
		logger.Initialize().Error(err)
	}
	logger.Initialize().Info("file created ", strg.path)

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	data := scanner.Bytes()

	err = json.Unmarshal(data, &strg.data)
	if err != nil {
		logger.Initialize().Info("Ошибка при unmarshal ", err)
	}
}
