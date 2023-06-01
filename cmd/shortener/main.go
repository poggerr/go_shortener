package main

import (
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/routers"
	"github.com/poggerr/go_shortener/internal/server"
)

func main() {
	cfg := config.NewConf()
	strg := storage.NewStorage(cfg.Path)

	r := routers.Router(cfg, strg)
	server.Server(cfg.Serv, r)
}
