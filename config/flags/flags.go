package flags

import (
	"flag"
	"os"
)

var mainFlags = flag.NewFlagSet("main", flag.ExitOnError)
var Serv string
var DefUrl string

func ParseFlags() {

	mainFlags.StringVar(&Serv, "a", ":8080", "write down server")
	mainFlags.StringVar(&DefUrl, "b", "http://localhost:8080/", "write down default url")
	err := mainFlags.Parse(os.Args[1:])
	if err != nil {
		return
	}
}
