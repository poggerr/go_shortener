package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/go_shortener/internal/app/shorten"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"io"
	"log"
	"net/http"
	"os"
)

type App struct {
	cfg     *config.Config
	storage storage.Storage
}

func NewApp(cfg *config.Config, strg storage.Storage) *App {
	return &App{
		cfg:     cfg,
		storage: strg,
	}
}

func (a *App) ReadOldUrl(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	ans := shorten.UnShoring(id, a.storage)
	res.Header().Set("Location", ans)
	res.WriteHeader(307)

}

func (a *App) CreateShortUrl(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		fmt.Println(err.Error())
		res.Write([]byte("Ошибка запроса"))
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	local := a.cfg.DefUrl()

	if local != "" {
		if string(local[len(local)-1]) != "/" {
			local += "/"
		}
	}

	short, err := shorten.Shorting(string(body), a.storage)

	if err != nil {
		http.Error(res, "Введите ссылку", http.StatusBadRequest)
		return
	}

	local += short

	res.Header().Set("content-type", "text/plain ")

	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(local))

}
