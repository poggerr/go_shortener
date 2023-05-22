package main

import (
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/logger"
	"github.com/poggerr/go_shortener/internal/routers"
	"github.com/poggerr/go_shortener/internal/server"
)

func main() {
	cfg := config.NewConf()
	strg := storage.NewStorage()
	sugar := logger.Initialize()

	r := routers.Router(cfg, strg, sugar)
	server.Server(cfg.Serv(), r, sugar)
}
