package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/poggerr/go_shortener/internal/authorization"
	"github.com/poggerr/go_shortener/internal/models"
	"github.com/poggerr/go_shortener/internal/utils"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

type URLShortener struct {
	linkRepo Repository
	baseURL  string
}

// NewURLShortener создает URLShortener и инициализирует его адресом, по которому будут доступны методы,
// и репозиторием хранения ссылок.
func NewURLShortener(base string, repo Repository) *URLShortener {
	hand := URLShortener{}
	hand.linkRepo = repo
	hand.baseURL = base

	return &hand
}

// CreateShortURL хендлер создания короткой ссылки
func (a *URLShortener) CreateShortURL(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	if !utils.IsURL(string(body)) {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Некорректный URL"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shortURL, err := a.linkRepo.Store(ctx, userID, string(body))
	if err != nil {
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte(a.baseURL + "/" + shortURL))
		return
	}

	res.Header().Set("content-type", "text/plain; charset=utf-8")

	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(a.baseURL + "/" + shortURL))

}

// DBConnect хендлер проверки подключения к БД
func (a *URLShortener) DBConnect(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := a.linkRepo.Ping(ctx); err != nil {
		log.Debug().Msg(fmt.Sprintf("database connection error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}

// CreateBatch Хендлер создания нескольких ссылок
func (a *URLShortener) CreateBatch(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())
	var list models.BatchList
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &list)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("unmarshal error: %s", err))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	list, err = a.linkRepo.StoreBatch(ctx, userID, list, a.baseURL)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshal, err := json.Marshal(list)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("marshal error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusCreated)
	res.Write(marshal)

}

// CreateJSONShorten хендлер создания ссылки из json
func (a *URLShortener) CreateJSONShorten(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	var url models.URL

	err = json.Unmarshal(body, &url)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("unmarshal error: %s", err))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if !utils.IsURL(url.LongURL) {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Некорректный URL"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shortURL, err := a.linkRepo.Store(ctx, userID, url.LongURL)
	switch err {
	case nil:
		res.WriteHeader(http.StatusCreated)
	case err:
		res.WriteHeader(http.StatusConflict)
	}

	shortURL = a.baseURL + "/" + shortURL
	shortenMap := make(map[string]string)
	shortenMap["result"] = shortURL

	marshal, err := json.Marshal(shortenMap)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("marshal error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json ")
	res.Write(marshal)

}

// GetUrlsByUser хендлер получения списка ссылок пользователя
func (a *URLShortener) GetUrlsByUser(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	strg, err := a.linkRepo.GetUserStorage(ctx, userID, a.baseURL)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("error: %s", err))
		res.WriteHeader(http.StatusNoContent)
		return
	}

	marshal, err := json.Marshal(strg)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("marshal error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusOK)
	res.Write(marshal)

}

// ReadOriginalURL хендлер получения оригинальной ссылки
func (a *URLShortener) ReadOriginalURL(res http.ResponseWriter, req *http.Request) {
	shortURL := chi.URLParam(req, "id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ans, err := a.linkRepo.Restore(ctx, shortURL)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("error: %s", err))
		res.WriteHeader(http.StatusNoContent)
		return
	}

	res.Header().Set("content-type", "text/plain ")
	res.Header().Set("Location", ans)
	res.WriteHeader(307)

}

// DeleteUrls хендлер удаления ссылок
func (a *URLShortener) DeleteUrls(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var keys []string

	err = json.Unmarshal(body, &keys)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("unmarshal error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.linkRepo.Delete(ctx, userID, keys)

	res.WriteHeader(http.StatusAccepted)
}

type Repository interface {
	// Store сохраняет оригинальную ссылку и возвращает id (токен) сокращенного варианта.
	Store(ctx context.Context, user *uuid.UUID, longURL string) (id string, err error)
	// Restore возвращает оригинальную ссылку по его id.
	Restore(ctx context.Context, id string) (link string, err error)
	// Delete - помечает ссылки удаленными.
	// Согласно заданию - результат работы пользователю не возвращается.
	Delete(ctx context.Context, user *uuid.UUID, ids []string)
	// GetUserStorage возвращает массив всех ранее сокращенных пользователей ссылок.
	GetUserStorage(ctx context.Context, user *uuid.UUID, defURL string) (map[string]string, error)
	// StoreBatch сохраняет пакет ссылок в хранилище и возвращает список пакет id.
	StoreBatch(ctx context.Context, user *uuid.UUID, batchIn models.BatchList, defURL string) (models.BatchList, error)
	// Ping проверяет готовность к работе репозитория.
	Ping(context.Context) error
	// Close завершает работу репозитория в стиле graceful shutdown.
	Close() error
}
