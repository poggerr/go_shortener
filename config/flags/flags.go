package flags

import (
	"flag"
	"os"
)

var mainFlags = flag.NewFlagSet("main", flag.ExitOnError)
var serv string
var defUrl string

func ParseFlags() {

	mainFlags.StringVar(&serv, "a", ":8080", "write down server")
	mainFlags.StringVar(&defUrl, "b", "http://localhost:8080", "write down default url")

	err := mainFlags.Parse(os.Args[1:])
	if err != nil {
		return
	}

	if servRunAddr := os.Getenv("SERVER_ADDRESS"); servRunAddr != "" {
		serv = servRunAddr
	}
	if baseRunAddr := os.Getenv("BASE_URL"); baseRunAddr != "" {
		defUrl = baseRunAddr
	}

	if string(defUrl[len(defUrl)-1]) != "/" {
		defUrl += "/"
	}
}

func ReturnServ() string {
	return serv
}

func ReturnDefUrl() string {
	return defUrl
}
