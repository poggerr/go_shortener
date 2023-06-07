package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/go_shortener/internal/app"
	"github.com/poggerr/go_shortener/internal/app/middlewares"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/gzip"
)

func Router(cfg *config.Config, strg *storage.Storage) chi.Router {
	r := chi.NewRouter()
	newApp := app.NewApp(cfg, strg)
	r.Use(middlewares.WithLogging, gzip.GzipMiddleware)
	r.Post("/", newApp.CreateShortURL)
	r.Post("/api/shorten", newApp.CreateJSONShorten)
	r.Get("/{id}", newApp.ReadOldURL)
	return r
}
