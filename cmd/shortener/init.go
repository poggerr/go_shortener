package main

import (
	"database/sql"
	"fmt"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/storage/database"
	"github.com/poggerr/go_shortener/internal/storage/file"
	"github.com/poggerr/go_shortener/internal/storage/memory"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	cfg = config.NewConf()
	initRepository()
}

func initRepository() {
	var (
		err error
		db  *sql.DB
	)

	cs := cfg.DB
	if len(cs) != 0 {
		db, err = sql.Open("pgx", cfg.DB)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	if db != nil {
		repo, err = database.NewStorage(db)
		if err == nil {
			log.Info().Msg("In file storage will be used")
			return
		}
	}

	filename := cfg.Path
	if len(filename) != 0 {
		repo, err = file.NewStorage(filename)
		if err == nil {
			log.Info().Msg("In file storage will be used")
			return
		}
	}

	repo = memory.NewStorage()
	log.Info().Msg("In memory storage will be used")
}
