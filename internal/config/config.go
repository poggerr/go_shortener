// Package config выполняет фунцию конфигурации приложения
package config

import (
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
)

// Config базовая структура конфигурации
type Config struct {
	Serv   string `env:"SERVER_ADDRESS"`
	DefURL string `env:"BASE_URL"`
	Path   string `env:"FILE_STORAGE_PATH"`
	DB     string `env:"DATABASE_DSN"`
}

func NewConf() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	var cfg Config

	flag.StringVar(&cfg.Serv, "a", ":8080", "write down server")
	flag.StringVar(&cfg.DefURL, "b", "http://localhost:8080", "write down default url")
	flag.StringVar(&cfg.Path, "f", "/tmp/short-url-db.json", "write down path to storage")
	flag.StringVar(&cfg.DB, "d", "", "write down db")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg
}
