package server

import (
	"github.com/poggerr/go_shortener/internal/logger"
	"log"
	"net/http"
)

func Server(addr string, hand http.Handler) {

	server := &http.Server{
		Addr:      addr,
		Handler:   hand,
		TLSConfig: nil,
		//IdleTimeout:       120 * time.Second,
		MaxHeaderBytes: 16 * 1024,
		//ReadHeaderTimeout: 10 * time.Second,
		//ReadTimeout:       10 * time.Second,
		//WriteTimeout:      10 * time.Second,
	}

	logger.Initialize().Info("Running server on: ", addr)

	log.Fatal(server.ListenAndServe())
}
