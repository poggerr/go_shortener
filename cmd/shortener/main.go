package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/poggerr/go_shortener/internal/app/service"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/routers"
	"github.com/poggerr/go_shortener/internal/server"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg := config.NewConf()

	if cfg.DB != "" {

		fmem, err := os.Create(`mem.pprof`)
		if err != nil {
			panic(err)
		}
		defer fmem.Close()
		runtime.GC() // получаем статистику по использованию памяти
		if err = pprof.WriteHeapProfile(fmem); err != nil {
			panic(err)
		}

		db, err := sql.Open("pgx", cfg.DB)
		if err != nil {
			fmt.Println(err.Error())
		}
		strg := storage.NewStorage(cfg.Path, db)

		repo := service.NewDeleter(strg)
		go repo.WorkerDeleteURLs()

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
