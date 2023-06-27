package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/routers"
	"github.com/poggerr/go_shortener/internal/server"
)

func main() {
	cfg := config.NewConf()

	if cfg.DB != "" {
		db, err := sql.Open("pgx", cfg.DB)
		if err != nil {
			fmt.Println(err.Error())
		}
		strg := storage.NewStorage(cfg.Path, db)

		strg.RestoreDB()

		if cfg.Path != "" {
			strg.RestoreFromFile()
		}
		r := routers.Router(cfg, strg, db)
		server.Server(cfg.Serv, r)
	} else {

		strg := storage.NewStorage(cfg.Path, nil)
		if cfg.Path != "" {
			strg.RestoreFromFile()
		}

		r := routers.Router(cfg, strg, nil)
		server.Server(cfg.Serv, r)
	}

}
