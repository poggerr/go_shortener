// Package config выполняет фунцию конфигурации приложения
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/rs/zerolog/log"
	"os"
)

// Config базовая структура конфигурации
type Config struct {
	Serv        string `env:"SERVER_ADDRESS" json:"server_address"`
	DefURL      string `env:"BASE_URL" json:"base_url"`
	Path        string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DB          string `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	ConfigJSON  string `env:"CONFIG"`
}

// NewConf конструктор конфигурации
func NewConf() *Config {
	cfg := new(Config)
	//if err := godotenv.Load(); err != nil {
	//	log.Print("No .env file found")
	//}

	flag.StringVar(&cfg.Serv, "a", ":8080", "write down server")
	flag.StringVar(&cfg.DefURL, "b", "http://localhost:8080", "write down default url")
	flag.StringVar(&cfg.Path, "f", "/tmp/short-url-db.json", "write down path to storage")
	flag.StringVar(&cfg.DB, "d", "host=localhost user=shortener password=password dbname=shortener sslmode=disable", "write down db")
	flag.BoolVar(&cfg.EnableHTTPS, "s", false, "write down enable https")
	flag.StringVar(&cfg.ConfigJSON, "c", "", "write down config json")
	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	if len(cfg.ConfigJSON) > 0 {
		cfg.GetJSONConfigData()
	}

	return cfg
}

func (c *Config) GetJSONConfigData() {
	log.Info().Msg("load data from json file")
	file, err := os.Open(c.ConfigJSON)
	if err != nil {
		log.Err(err).Msgf("can't open file: %s", c.ConfigJSON)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}(file)

	err = json.NewDecoder(file).Decode(c)
	if err != nil {
		log.Err(err).Msgf("can't read config from given file: %s", c.ConfigJSON)
	}
}
