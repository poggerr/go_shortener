package config

import (
	"flag"
	"os"
)

var serv string
var defUrl string

func ParseServ() {

	flag.StringVar(&serv, "a", ":8080", "write down server")
	flag.StringVar(&defUrl, "b", "http://localhost:8080", "write down default url")
	flag.Parse()

	if servRunAddr := os.Getenv("SERVER_ADDRESS"); servRunAddr != "" {
		serv = servRunAddr
	}
	if baseRunAddr := os.Getenv("BASE_URL"); baseRunAddr != "" {
		defUrl = baseRunAddr
	}
}

func GetServ() string {
	return serv
}

func GetDefUrl() string {
	return defUrl
}
