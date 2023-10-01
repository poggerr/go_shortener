package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/poggerr/go_shortener/internal/app/models"
	"github.com/poggerr/go_shortener/internal/app/service"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/authorization"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/logger"
)

type App struct {
	cfg     *config.Config
	storage *storage.Storage
	db      *sql.DB
	repo    *service.URLRepo
}

func NewApp(cfg *config.Config, strg *storage.Storage, db *sql.DB, repo *service.URLRepo) *App {
	return &App{
		cfg:     cfg,
		storage: strg,
		db:      db,
		repo:    repo,
	}
}

func (a *App) ReadOldURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	ans, isDelete, err := service.Check(id, a.storage)
	if isDelete {
		res.Header().Set("content-type", "text/plain ")
		res.WriteHeader(http.StatusGone)
		return
	}
	if err != nil {
		fmt.Fprint(res, err.Error())
		logger.Initialize().Info(err)
		return
	}

	res.Header().Set("content-type", "text/plain ")

	res.Header().Set("Location", ans)
	res.WriteHeader(307)

}

func (a *App) CreateShortURL(res http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session_token")
	var userID string
	if err != nil {
		logger.Initialize().Info(err)
	}
	switch c {
	case nil:
		uuidUserID := uuid.New()
		jwtString, err := authorization.BuildJWTString(&uuidUserID)
		if err != nil {
			logger.Initialize().Info(err)
		}

		cook := &http.Cookie{
			Name:    "session_token",
			Value:   jwtString,
			Path:    "/",
			Domain:  "localhost",
			Expires: time.Now().Add(120 * time.Second),
		}

		http.SetCookie(res, cook)
		userID = uuidUserID.String()
	default:
		userID = authorization.GetUserID(c.Value)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	short, err := service.CreateService(string(body), a.storage, userID)
	shortURL := a.cfg.DefURL + "/" + short
	if err != nil {
		logger.Initialize().Info(err)
		res.Header().Set("content-type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte(shortURL))
		return
	}

	res.Header().Set("content-type", "text/plain; charset=utf-8")

	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(shortURL))

}

func (a *App) CreateJSONShorten(res http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session_token")
	var userID string
	if err != nil {
		logger.Initialize().Info(err)
	}
	switch c {
	case nil:
		uuidUserID := uuid.New()
		jwtString, err := authorization.BuildJWTString(&uuidUserID)
		if err != nil {
			logger.Initialize().Info(err)
		}

		cook := &http.Cookie{
			Name:    "session_token",
			Value:   jwtString,
			Path:    "/",
			Domain:  "localhost",
			Expires: time.Now().Add(120 * time.Second),
		}

		http.SetCookie(res, cook)
		userID = uuidUserID.String()
	default:
		userID = authorization.GetUserID(c.Value)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	var url models.URL

	err = json.Unmarshal(body, &url)
	if err != nil {
		logger.Initialize().Info(err)
	}

	short, err := service.CreateService(url.LongURL, a.storage, userID)
	shortURL := a.cfg.DefURL + "/" + short

	if err != nil {
		shortenMap := make(map[string]string)

		shortenMap["result"] = shortURL

		marshal, err := json.Marshal(shortenMap)
		if err != nil {
			logger.Initialize().Info(err)
		}

		res.Header().Set("content-type", "application/json ")
		res.WriteHeader(http.StatusConflict)
		res.Write(marshal)
		return
	}
	shortenMap := make(map[string]string)

	shortenMap["result"] = shortURL

	marshal, err := json.Marshal(shortenMap)
	if err != nil {
		logger.Initialize().Info(err)
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusCreated)

	res.Write(marshal)

}

func (a *App) DBConnect(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := a.db.PingContext(ctx); err != nil {
		logger.Initialize().Error("Ошибка при подключении к БД ", err)
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}

func (a *App) CreateBatch(res http.ResponseWriter, req *http.Request) {
	var list models.BatchList
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &list)
	if err != nil {
		logger.Initialize().Info(err)
	}

	list = service.SaveMultipleToDB(list, a.storage, a.cfg.DefURL)

	marshal, err := json.Marshal(list)
	if err != nil {
		logger.Initialize().Info("Ошибка при формировании ответа ", err)
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusCreated)
	res.Write(marshal)

}

func (a *App) GetUrlsByUser(res http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session_token")
	var userID string
	if err != nil {
		logger.Initialize().Info(err)
		res.WriteHeader(http.StatusNoContent)
	}
	if c != nil {
		userID = authorization.GetUserID(c.Value)
	}
	if userID == "" {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("Пользователь не авторизован!"))
		return
	}

	strg := a.storage.GetUrlsByUsesID(userID, a.cfg.DefURL)

	marshal, err := json.Marshal(strg)
	if err != nil {
		logger.Initialize().Info(err)
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusOK)
	res.Write(marshal)

}

func (a *App) DeleteUrls(res http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session_token")
	var userID string
	if err != nil {
		logger.Initialize().Info(err)
		res.WriteHeader(http.StatusNoContent)
	}
	if c != nil {
		userID = authorization.GetUserID(c.Value)
	}
	if userID == "" {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("Пользователь не авторизован!"))
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	var keys []string

	err = json.Unmarshal(body, &keys)
	if err != nil {
		logger.Initialize().Info(err)
	}

	err = a.repo.DeleteAsync(keys, userID)
	if err != nil {
		logger.Initialize().Info(err)
	}

	res.WriteHeader(http.StatusAccepted)
}
