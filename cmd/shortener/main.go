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

	fmt.Println("Running server on", flags.Serv)
	fmt.Println("DefUrl: ", flags.DefUrl)

	log.Fatal(http.ListenAndServe(flags.Serv, r))
}
