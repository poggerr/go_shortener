package app

import (
	"database/sql"
	"github.com/poggerr/go_shortener/internal/async"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/storage"
)

type App struct {
	cfg     *config.Config
	storage *storage.Storage
	db      *sql.DB
	repo    *async.URLRepo
}

func NewApp(cfg *config.Config, strg *storage.Storage, db *sql.DB, repo *async.URLRepo) *App {
	return &App{
		cfg:     cfg,
		storage: strg,
		db:      db,
		repo:    repo,
	}
}
