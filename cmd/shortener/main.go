package main

import (
	"fmt"
	"github.com/poggerr/go_shortener/config/flags"
	"github.com/poggerr/go_shortener/internal/routers"
	"log"
	"net/http"
)

func main() {
	flags.ParseFlags()

	r := routers.Routers()

	fmt.Println("Running server on", flags.ReturnServ())

	log.Fatal(http.ListenAndServe(flags.ReturnServ(), r))
}
