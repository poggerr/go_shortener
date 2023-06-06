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

	if cfg.Serv == "" {
		flag.StringVar(&cfg.Serv, "a", ":8080", "write down server")
	}
	if cfg.DefUrl == "" {
		flag.StringVar(&cfg.DefUrl, "b", "http://localhost:8080", "write down default url")
	}
	if cfg.Path == "" {
		flag.StringVar(&cfg.Path, "f", "", "write down path to storage")
		fmt.Println(cfg.Path)
	}
	flag.Parse()

	return &cfg
}
