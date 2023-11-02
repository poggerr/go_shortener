package main

import (
	"context"
	"github.com/poggerr/go_shortener/internal/server"
	"github.com/poggerr/go_shortener/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/poggerr/go_shortener/internal/config"
)

var (
	cfg  *config.Config
	repo service.URLShortenerService
)

func main() {
	srv := CreateServer()
	Run(srv)
}

func CreateServer() *http.Server {
	return server.Server(cfg.Serv, cfg.DefURL, repo)
}

func Run(srv *http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("listen: %+v\n", err)
		}
		log.Info().Msg("Server started")
	}()

	<-ctx.Done()

	log.Info().Msg("Server stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		err := repo.Close()
		if err != nil {
			log.Error().Msgf("Caught an error due closing repository:%+v", err)
		}

		log.Info().Msg("Everything is closed properly")
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Msgf("Server Shutdown Failed:%+v", err)
	}
	stop()
	log.Info().Msg("Server exited properly")
}
