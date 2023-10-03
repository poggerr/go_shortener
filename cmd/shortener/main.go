package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/poggerr/go_shortener/internal/async"
	"github.com/poggerr/go_shortener/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
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

		repo := async.NewDeleter(strg)

		baseCTX := context.Background()
		ctx, cancelFunction := context.WithCancel(baseCTX)
		defer func() {
			cancelFunction()
		}()

		go repo.WorkerDeleteURLs(ctx)

		strg.RestoreDB()

		if cfg.Path != "" {
			strg.RestoreFromFile()
		}
		r := routers.Router(cfg, strg, db, repo)
		server.Server(cfg.Serv, r)

	} else {

		strg := storage.NewStorage(cfg.Path, nil)
		if cfg.Path != "" {
			strg.RestoreFromFile()
		}

		r := routers.Router(cfg, strg, nil, nil)
		server.Server(cfg.Serv, r)
	}

}
