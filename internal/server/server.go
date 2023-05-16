package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Server(addr string, hand http.Handler) {

	server := &http.Server{
		Addr:                         addr,
		Handler:                      hand,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		IdleTimeout:                  120 * time.Second,
		MaxHeaderBytes:               16 * 1024,
		ReadHeaderTimeout:            10 * time.Second,
		ReadTimeout:                  10 * time.Second,
		WriteTimeout:                 10 * time.Second,
	}

	fmt.Println("Running server on", addr)

	log.Fatal(server.ListenAndServe())
}
