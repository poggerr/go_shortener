package storage

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/poggerr/go_shortener/internal/logger"
	"os"
	"path"
	"time"
)

type Storage struct {
	data map[string]string
	path string
	DB   *sql.DB
}

func NewStorage(p string, db *sql.DB) *Storage {
	return &Storage{
		data: make(map[string]string),
		path: p,
		DB:   db,
	}
}

func (strg *Storage) Save(key, value string) {
	strg.data[key] = value
	if strg.path != "" {
		strg.SaveToFile()
	}
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

func (strg *Storage) RestoreDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := strg.DB.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS urls (
	    "correlation_id" TEXT,
		"longurl" TEXT,
		"shorturl" TEXT
	)
	`)
	if err != nil {
		logger.Initialize().Info("Ошибка при создании таблицы urls ", err)
	}

}

func (strg *Storage) SaveToDB(longurl, shorturl string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf("INSERT INTO urls (longurl, shorturl) VALUES ('%s', '%s')", longurl, shorturl)

	_, err := strg.DB.ExecContext(ctx, query)
	if err != nil {
		logger.Initialize().Info("Ошибка при записи в urls ", err)
	}
}
