package server

import (
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func Server(addr string, hand http.Handler, sugar *zap.SugaredLogger) {

	server := &http.Server{
		Addr:              addr,
		Handler:           hand,
		TLSConfig:         nil,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    16 * 1024,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	sugar.Info("Running server on: ", addr)

	log.Fatal(server.ListenAndServe())
}
