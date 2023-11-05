package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/poggerr/go_shortener/internal/models"
	"github.com/poggerr/go_shortener/internal/service"
	"github.com/poggerr/go_shortener/internal/utils"
	"time"
)

// TODO переписать
type Storage struct {
	database *sql.DB
	done     chan bool
	delBatch chan userID
}

var _ service.URLShortenerService = (*Storage)(nil)

func NewStorage(db *sql.DB) (*Storage, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	tx.Commit()
	return &Storage{
		database: db,
	}, nil
}

func (strg *Storage) Store(ctx context.Context, user *uuid.UUID, longURL string) (id string, err error) {
	shortURL := utils.NewShortURL()

	_, err = strg.database.ExecContext(ctx, "INSERT INTO urls (long_url, short_url, user_id) VALUES ($1, $2, $3)", longURL, shortURL, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			ans := strg.database.QueryRowContext(ctx, "SELECT short_url FROM urls WHERE long_url=$1", longURL)
			errScan := ans.Scan(&shortURL)
			if errScan != nil {
				fmt.Println(err)
			}
			return shortURL, err
		}
	}
	return shortURL, err
}

func (strg *Storage) Restore(ctx context.Context, shortURL string) (link string, err error) {
	var isDelete bool
	ans := strg.database.QueryRowContext(ctx, "SELECT long_url, is_deleted FROM urls WHERE short_url=$1", shortURL)
	errScan := ans.Scan(&link, &isDelete)
	switch {
	case errScan != nil:
		return "", errScan
	case isDelete:
		return "", errors.New("ссылка удалена")
	}
	return
}

func (strg *Storage) Delete(ctx context.Context, user *uuid.UUID, ids []string) {
	return
}

//// UserURLs структура для удаления ссылок
//type UserURLs struct {
//	UserID *uuid.UUID
//	URLs   []string
//}
//
//// DeleteUrls удаление ссылок
//
//// TODO переписать удаление
//func (strg *Storage) DeleteUrls(mas UserURLs) {
//	tx, err := strg.DB.Begin()
//	if err != nil {
//		logger.Initialize().Error(err)
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	for _, m := range mas.URLs {
//		_, err = tx.ExecContext(ctx, "UPDATE urls SET is_deleted=true WHERE short_url=$1 AND user_id=$2", m, mas.UserID)
//		if err != nil {
//			logger.Initialize().Info("Ошибка при удалении", err)
//		}
//	}
//	tx.Commit()
//}

func (strg *Storage) GetUserStorage(ctx context.Context, user *uuid.UUID, defURL string) (map[string]string, error) {
	rows, err := strg.database.QueryContext(ctx, "SELECT * FROM urls WHERE user_id=$1", user)
	if err != nil {
		return nil, err
	}

	userStorage := make(map[string]string)
	for rows.Next() {
		var url models.Urls
		if err = rows.Scan(&url.UserID, &url.LongURL, &url.ShortURL, &url.DeletedFlag); err != nil {
			return nil, err
		}
		userStorage[url.ShortURL] = url.LongURL
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(userStorage) < 1 {
		return nil, errors.New("у пользователя нет ссылок")
	}

	return userStorage, nil
}

func (strg *Storage) StoreBatch(ctx context.Context, user *uuid.UUID, batchIn models.BatchList, defURL string) (models.BatchList, error) {
	tx, err := strg.database.Begin()
	if err != nil {
		return nil, err
	}
	for i, v := range batchIn {
		shortURL, err := strg.Store(ctx, user, v.OriginalURL)
		if err != nil {
			return nil, err
		}
		batchIn[i].ShortURL = defURL + "/" + shortURL
	}
	tx.Commit()
	return batchIn, nil
}

func (strg *Storage) Ping(ctx context.Context) error {
	return strg.database.PingContext(ctx)
}

type userID struct {
	User string
	ID   string
}

// Close закрывает базу данных
// TODO переписать
func (s *Storage) Close() error {
	s.done <- true
	// важен порядок закрытия!
	close(s.delBatch)
	close(s.done)
	return s.database.Close()
}

//type URLRepo struct {
//	urlsToDeleteChan chan storage.UserURLs
//	repository       storage.Storage
//}
//
//func NewDeleter(strg *storage.Storage) *URLRepo {
//	return &URLRepo{
//		urlsToDeleteChan: make(chan storage.UserURLs, 10),
//		repository:       *strg,
//	}
//}
//
//func (r *URLRepo) DeleteAsync(ids []string, userID *uuid.UUID) error {
//	r.urlsToDeleteChan <- storage.UserURLs{UserID: userID, URLs: ids}
//	return nil
//}
//
//// WorkerDeleteURLs воркер удаления ссылок
//func (r *URLRepo) WorkerDeleteURLs(ctx context.Context) {
//	for urls := range r.urlsToDeleteChan {
//		select {
//		case <-ctx.Done():
//			return
//		default:
//			r.repository.DeleteUrls(urls)
//		}
//	}
//}
