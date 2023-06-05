package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v8"
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

	if cfg.Serv == "" || cfg.DefUrl == "" || cfg.Path == "" {
		flag.StringVar(&cfg.Serv, "a", ":8080", "write down server")
		flag.StringVar(&cfg.Path, "f", "/tmp/short-url-db.json", "write down path to storage")
		flag.StringVar(&cfg.DefUrl, "b", "http://localhost:8080", "write down default url")
		flag.Parse()
	}

	return &cfg
}
