package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v8"
)

type Config struct {
	serv   string `env:"SERVER_ADDRESS"`
	defUrl string `env:"BASE_URL"`
}

func (cfg Config) Add() Config {
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	flag.StringVar(&cfg.serv, "a", ":8080", "write down server")
	flag.StringVar(&cfg.defUrl, "b", "http://localhost:8080", "write down default url")
	flag.Parse()
	return cfg
}

func (cfg Config) Serv() string {
	return cfg.serv
}

func (cfg Config) DefUrl() string {
	return cfg.defUrl
}

func NewCfg() Config {
	var cfg Config
	cfg = cfg.Add()
	return cfg
}

func NewCfgForTests() Config {
	return Config{
		serv:   ":8080",
		defUrl: "http://localhost:8080",
	}
}
