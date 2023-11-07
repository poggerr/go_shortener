package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/poggerr/go_shortener/internal/handlers"
	"github.com/poggerr/go_shortener/internal/models"
	"github.com/poggerr/go_shortener/internal/service"
	"github.com/poggerr/go_shortener/internal/utils"
	"github.com/rs/zerolog/log"
	"time"
)

type Storage struct {
	database *sql.DB
	done     chan bool
	delBatch chan userID
}

var _ service.URLShortenerService = (*Storage)(nil)

func NewStorage(db *sql.DB) (*Storage, error) {
	s := Storage{database: db}
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

	s.delBatch = make(chan userID)
	s.done = make(chan bool)
	go s.deleteConsume()

	return &s, nil
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
		return "", handlers.ErrLinkIsDeleted
	}
	return
}

func (strg *Storage) Delete(_ context.Context, user *uuid.UUID, ids []string) {
	ch := make(chan userID)
	go strg.deleteProduce(user, ids, ch)
	go strg.deleteWork(ch)

}

func (strg *Storage) deleteProduce(user *uuid.UUID, ids []string, ch chan userID) {
	for i, id := range ids {
		ch <- userID{User: user, ID: id}
		log.Debug().Msgf("%v", i)
	}
	close(ch)
}

func (strg *Storage) deleteWork(ch chan userID) {
	for id := range ch {
		strg.delBatch <- id
	}
}

func (strg *Storage) deleteConsume() {
	flush := func() {
		for {
			time.Sleep(time.Second)
			strg.done <- true
		}
	}

	go flush()

	var buf = make([]userID, 10)
	i := 0
	for {
		select {
		case <-strg.done:
			if i != 0 {
				log.Debug().Msg(fmt.Sprint(buf[:i]))
				err := strg.deleteBatch(buf[:i])
				if err != nil {
					log.Err(err).Send()
				}
				i = 0
			}
		case id, ok := <-strg.delBatch:
			if !ok {
				return
			}
			if i == len(buf) {
				log.Debug().Msg(fmt.Sprint(buf))
				err := strg.deleteBatch(buf)
				if err != nil {
					log.Err(err).Send()
				}
				i = 0
			}
			buf[i] = id
			i++
		}
	}
}

func (strg *Storage) deleteBatch(ids []userID) error {
	// шаг 1 — объявляем транзакцию
	tx, err := strg.database.Begin()
	if err != nil {
		return err
	}

	// Это чтобы мы тут тоже не зависли надолго
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	// шаг 2 — готовим инструкцию
	stmt, err := tx.PrepareContext(ctx, "UPDATE urls SET is_deleted=true WHERE short_url=$1 AND user_id=$2")
	if err != nil {
		return err
	}
	// шаг 2.1 — не забываем закрыть инструкцию, когда она больше не нужна
	defer func() {
		if err = stmt.Close(); err != nil {
			log.Err(err).Send()
		}
	}()

	// шаг 3 - выполняем задачу
	for _, id := range ids {
		_, err = stmt.ExecContext(ctx, id.ID, id.User)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// шаг 4 — сохраняем изменения
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (strg *Storage) GetUserStorage(ctx context.Context, user *uuid.UUID, _ string) (map[string]string, error) {
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
	User *uuid.UUID
	ID   string
}

// Close закрывает базу данных
func (strg *Storage) Close() error {
	strg.done <- true
	// важен порядок закрытия!
	close(strg.delBatch)
	fmt.Println("delBatch close")
	close(strg.done)
	fmt.Println("done close")
	return strg.database.Close()
}
