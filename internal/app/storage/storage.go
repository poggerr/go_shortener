package storage

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/poggerr/go_shortener/internal/logger"
	"os"
	"path"
	"time"
)

type URL struct {
	LongURL  string `json:"longUrl"`
	ShortURL string `json:"shortUrl"`
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

func (strg *Storage) LongURL(key string) (string, error) {
	val, ok := strg.data[key]
	if !ok {
		return "/", errors.New("такой ссылки нет. Введите запрос повторно")
	}
	return val, nil
}

func (strg *Storage) SaveToFile() {
	file, err := os.OpenFile(strg.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Initialize().Error(err)
		}
	}(file)
	if err != nil {
		logger.Initialize().Error(err)
	}

	data, _ := json.Marshal(strg.data)

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

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	data := scanner.Bytes()

	err = json.Unmarshal(data, &strg.data)
	if err != nil {
		logger.Initialize().Info("Ошибка при unmarshal ", err)
	}
}

func RestoreDB(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS urls (
		"id" INTEGER PRIMARY KEY,
		"longurl" VARCHAR(250) NOT NULL DEFAULT '',
		"shorturl" VARCHAR(250) NOT NULL DEFAULT ''
	)
	`)
	if err != nil {
		logger.Initialize().Info("Ошибка при создании таблицы urls ", err)
	}

}
