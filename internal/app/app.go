package app

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/go_shortener/internal/app/shorten"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/logger"
	"io"
	"net/http"
)

type App struct {
	cfg     *config.Config
	storage *storage.Storage
}

func NewApp(cfg *config.Config, strg *storage.Storage) *App {
	return &App{
		cfg:     cfg,
		storage: strg,
	}
}

func (a *App) ReadOldURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	ans, err := shorten.UnShoring(id, a.storage)
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

	short := shorten.Shorting(string(body), a.storage)

	res.Header().Set("content-type", "text/plain; charset=utf-8")

	res.WriteHeader(http.StatusCreated)

	fmt.Fprint(res, a.cfg.DefURL, "/", short)

}

type URL struct {
	LongURL  string `json:"url"`
	ShortURL string `json:"result"`
}

func (a *App) CreateJSONShorten(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	var url URL

	err = json.Unmarshal(body, &url)
	if err != nil {
		logger.Initialize().Info(err)
	}

	shortURL := shorten.Shorting(url.LongURL, a.storage)
	shortenMap := make(map[string]string)

	shortURL = a.cfg.DefURL + "/" + shortURL

	shortenMap["result"] = shortURL

	marshal, err := json.Marshal(shortenMap)
	if err != nil {
		logger.Initialize().Info(err)
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusCreated)

	res.Write(marshal)

}
