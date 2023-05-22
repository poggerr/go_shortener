package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/go_shortener/internal/app"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/logger"
	"go.uber.org/zap"
)

func Router(cfg *config.Config, strg *storage.Storage, sugar *zap.SugaredLogger) chi.Router {
	r := chi.NewRouter()
	newApp := app.NewApp(cfg, strg)
	r.Use(logger.WithLoggingReq)
	r.Use(logger.WithLoggingRes)
	r.Route("/", func(r chi.Router) {
		r.Post("/", newApp.CreateShortUrl)
		r.Get("/{id}", newApp.ReadOldUrl)
	})
	return r
}
