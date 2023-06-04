package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/go_shortener/internal/app/shorten"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/logger"
	"io"
	"log"
	"net/http"
	"os"
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

func (a *App) ReadOldUrl(res http.ResponseWriter, req *http.Request) {
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

func (a *App) CreateShortUrl(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Fprint(res, err.Error())
		logger.Initialize().Info(err)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	short := shorten.Shorting(string(body), a.storage)

	res.Header().Set("content-type", "text/plain; charset=utf-8")

	res.WriteHeader(http.StatusCreated)

	fmt.Fprint(res, a.cfg.DefUrl, "/", short)

}

func (a *App) CreateJsonShorten(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		fmt.Println(err.Error())
		res.Write([]byte(err.Error()))
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	short, err := shorten.JsonCreater(body, a.storage, a.cfg.DefUrl)
	if err != nil {
		return
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusCreated)

	res.Write(short)

}
