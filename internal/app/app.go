package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/poggerr/go_shortener/internal/app/models"
	"github.com/poggerr/go_shortener/internal/app/service"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/authorization"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/logger"
	"io"
	"net/http"
	"time"
)

type App struct {
	cfg     *config.Config
	storage *storage.Storage
	db      *sql.DB
}

func NewApp(cfg *config.Config, strg *storage.Storage, db *sql.DB) *App {
	return &App{
		cfg:     cfg,
		storage: strg,
		db:      db,
	}
}

func (a *App) ReadOldURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	ans, err := service.Take(id, a.storage)
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
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	userId := uuid.New()

	short, err := service.ServiceCreate(string(body), a.storage, userId.String())
	shortURL := a.cfg.DefURL + "/" + short
	if err != nil {
		logger.Initialize().Info(err)
		res.Header().Set("content-type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte(shortURL))
		return
	}

	jwtString, err := authorization.BuildJWTString(&userId)
	if err != nil {
		logger.Initialize().Info(err)
	}

	c := &http.Cookie{
		Name:    "session_token",
		Value:   jwtString,
		Path:    "/",
		Domain:  "localhost",
		Expires: time.Now().Add(120 * time.Second),
	}

	http.SetCookie(res, c)

	res.Header().Set("content-type", "text/plain; charset=utf-8")

	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(shortURL))

}

func (a *App) CreateJSONShorten(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	userId := uuid.New()

	var url models.URL

	err = json.Unmarshal(body, &url)
	if err != nil {
		logger.Initialize().Info(err)
	}

	short, err := service.ServiceCreate(url.LongURL, a.storage, userId.String())
	shortURL := a.cfg.DefURL + "/" + short

	if err != nil {
		shortenMap := make(map[string]string)

		shortenMap["result"] = shortURL

		marshal, err := json.Marshal(shortenMap)
		if err != nil {
			logger.Initialize().Info(err)
		}

		jwtString, err := authorization.BuildJWTString(&userId)
		if err != nil {
			logger.Initialize().Info(err)
		}

		c := &http.Cookie{
			Name:    "session_token",
			Value:   jwtString,
			Path:    "/",
			Domain:  "localhost",
			Expires: time.Now().Add(120 * time.Second),
		}

		http.SetCookie(res, c)

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

	jwtString, err := authorization.BuildJWTString(&userId)
	if err != nil {
		logger.Initialize().Info(err)
	}

	c := &http.Cookie{
		Name:    "session_token",
		Value:   jwtString,
		Path:    "/",
		Domain:  "localhost",
		Expires: time.Now().Add(120 * time.Second),
	}

	http.SetCookie(res, c)

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

	list = service.SaveMultipleToDB(list, a.storage)

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
	var userId string
	if err != nil {
		logger.Initialize().Info(err)
		res.WriteHeader(http.StatusNoContent)
	}
	if c != nil {
		userId = authorization.GetUserID(c.Value)
	}
	if userId == "" {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("Пользователь не авторизован!"))
		return
	}

	strg := a.storage.GetUrlsByUsesId(userId)

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
	var userId string
	if err != nil {
		logger.Initialize().Info(err)
		res.WriteHeader(http.StatusNoContent)
	}
	if c != nil {
		userId = authorization.GetUserID(c.Value)
	}
	if userId == "" {
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

	service.ServiceDelete(keys, userId, a.storage)

	res.WriteHeader(http.StatusAccepted)
}
