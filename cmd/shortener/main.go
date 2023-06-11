package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/logger"
	"github.com/poggerr/go_shortener/internal/routers"
	"github.com/poggerr/go_shortener/internal/server"
)

func main() {
	cfg := config.NewConf()

	db, err := sql.Open("pgx", cfg.DB)
	if err != nil {
		logger.Initialize().Error("Ошибка при подключении к БД ", err)
	}
	defer db.Close()

	storage.RestoreDB(db)

	strg := storage.NewStorage(cfg.Path)
	if cfg.Path != "" {
		strg.RestoreFromFile()
	}

	r := routers.Router(cfg, strg, db)
	server.Server(cfg.Serv, r)
}
