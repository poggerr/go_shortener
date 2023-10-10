package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/poggerr/go_shortener/internal/logger"
)

// RestoreFromFile восстанавливает данные из файла
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

// RestoreDB восстанавливает базу данных
func (strg *Storage) RestoreDB() {
	tx, err := strg.DB.Begin()
	if err != nil {
		logger.Initialize().Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS users (
	    "id" UUID UNIQUE,
	    "username" TEXT,
		"pass" TEXT
	)
	`)
	if err != nil {
		logger.Initialize().Info("Ошибка при создании таблицы users ", err)
	}

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS urls (
	    "user_id" UUID,
		"long_url" TEXT UNIQUE,
		"short_url" TEXT,
		"is_deleted" BOOL DEFAULT false
	)
	`)
	if err != nil {
		logger.Initialize().Info("Ошибка при создании таблицы urls ", err)
	}

	tx.Commit()

}
