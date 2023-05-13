package main

import (
	"fmt"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/routers"
	"log"
	"net/http"
)

func main() {
	config.ParseServ()
	serv := config.GetServ()

	r := routers.Routers()

	fmt.Println("Running server on", serv)

	log.Fatal(http.ListenAndServe(serv, r))
}
