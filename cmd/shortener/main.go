package main

import (
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/routers"
	"github.com/poggerr/go_shortener/internal/server"
)

func main() {
	cfg := config.NewCfg()
	strg := storage.NewStorage()

	r := routers.Router(&cfg, strg)
	server.Server(cfg.Serv(), r)
}
