package storage

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/poggerr/go_shortener/internal/app/models"
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
	    "user_id" TEXT,
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

func (strg *Storage) SaveToDB(longurl, shorturl string, userID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := strg.DB.ExecContext(ctx, "INSERT INTO urls (long_url, short_url, user_id) VALUES ($1, $2, $3)", longurl, shorturl, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			ans := strg.DB.QueryRowContext(ctx, "SELECT short_url FROM urls WHERE long_url=$1", longurl)
			errScan := ans.Scan(&shorturl)
			if errScan != nil {
				logger.Initialize().Info(errScan)
			}
			return shorturl, err
		}
	}
	return "", err
}

func (strg *Storage) CreateUser(username, pass string, id *uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := strg.DB.ExecContext(ctx, "INSERT INTO users (id, username, pass) VALUES ($1, $2, $3)", id, username, pass)
	if err != nil {
		logger.Initialize().Info("Ошибка при создании юзера ", err)
		return err
	}
	return nil
}

func (strg *Storage) GetUserID(username string) *uuid.UUID {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id *uuid.UUID
	ans := strg.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE username=$1", username)
	errScan := ans.Scan(&id)
	if errScan != nil {
		logger.Initialize().Info(errScan)
	}
	return id
}

func (strg *Storage) TakeLongURLIsDelete(shortURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var isDelete bool
	ans := strg.DB.QueryRowContext(ctx, "SELECT is_deleted FROM urls WHERE short_url=$1", shortURL)
	errScan := ans.Scan(&isDelete)
	if errScan != nil {
		logger.Initialize().Info(errScan)
	}
	return isDelete
}

func (strg *Storage) GetUrlsByUsesID(id string, defURL string) *models.Storage {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := strg.DB.QueryContext(ctx, "SELECT * FROM urls WHERE user_id=$1", id)
	if err != nil {
		logger.Initialize().Info(err)
	}

	storage := make(models.Storage, 0)
	for rows.Next() {
		var url models.Urls
		if err = rows.Scan(&url.UserID, &url.LongURL, &url.ShortURL, &url.DeletedFlag); err != nil {
			logger.Initialize().Info(err)
		}
		url.ShortURL = defURL + "/" + url.ShortURL
		storage = append(storage, url)
	}

	if err = rows.Err(); err != nil {
		logger.Initialize().Info(err)
	}

	return &storage
}

type UserURLs struct {
	UserID string
	URLs   []string
}

func (strg *Storage) DeleteUrls(mas UserURLs) {
	tx, err := strg.DB.Begin()
	if err != nil {
		logger.Initialize().Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, m := range mas.URLs {
		_, err = tx.ExecContext(ctx, "UPDATE urls SET is_deleted=true WHERE short_url=$1 AND user_id=$2", m, mas.UserID)
		if err != nil {
			logger.Initialize().Info("Ошибка при удалении", err)
		}
	}

	tx.Commit()

}
