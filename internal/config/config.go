package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v8"
	"os"
	"path"
	"strings"
)

type Config struct {
	serv   string `env:"SERVER_ADDRESS"`
	defUrl string `env:"BASE_URL"`
	path   string `env:"FILE_STORAGE_PATH"`
}

func NewConf() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	flag.StringVar(&cfg.serv, "a", ":8080", "write down server")
	flag.StringVar(&cfg.defUrl, "b", "http://localhost:8080", "write down default url")
	flag.StringVar(&cfg.path, "f", "/tmp/short-url-db.json", "write down path to storage")
	flag.Parse()

	dir, _ := path.Split(cfg.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0666)
		if err != nil {
			fmt.Printf(err.Error())
		}
	}

	if !strings.Contains(cfg.path, ".json") {
		cfg.path += ".json"
	}

	return &cfg
}

func (cfg Config) Serv() string {
	return cfg.serv
}

func (cfg Config) DefUrl() string {
	return cfg.defUrl
}

func (cfg Config) Path() string {
	return cfg.path
}

func NewDefConf() Config {
	return Config{
		serv:   ":8080",
		defUrl: "http://localhost:8080",
		path:   "/tmp/short-url-db3.json",
	}
}
