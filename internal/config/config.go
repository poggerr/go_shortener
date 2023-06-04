package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/poggerr/go_shortener/internal/logger"
	"os"
	"path"
)

type Config struct {
	Serv   string `env:"SERVER_ADDRESS"`
	DefUrl string `env:"BASE_URL"`
	Path   string `env:"FILE_STORAGE_PATH"`
}

func NewConf() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	flag.StringVar(&cfg.Serv, "a", ":8080", "write down server")
	flag.StringVar(&cfg.DefUrl, "b", "http://localhost:8080", "write down default url")
	flag.StringVar(&cfg.Path, "f", "/tmp/short-url-db.json", "write down path to storage")
	flag.Parse()

	dir, _ := path.Split(cfg.Path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0666)
		if err != nil {
			logger.Initialize().Info(err)
		}
	}

	return &cfg
}
