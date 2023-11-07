package main

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/server"
	"github.com/poggerr/go_shortener/internal/service"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
	cfg          *config.Config
	repo         service.URLShortenerService
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	srv := CreateServer()
	Run(srv)
}

func CreateServer() *http.Server {
	return server.Server(cfg.Serv, cfg.DefURL, repo)
}

func Run(srv *http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return srv.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return srv.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		repo.Close()
		fmt.Printf("exit reason: %s \n", err)
	}

	//go func() {
	//	err := srv.ListenAndServe()
	//	if err != nil && err != http.ErrServerClosed {
	//		log.Fatal().Msgf("listen: %+v\n", err)
	//	}
	//	log.Info().Msg("Server started")
	//}()
	//
	//<-ctx.Done()
	//
	//log.Info().Msg("Server stopped")
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer func() {
	//	err := repo.Close()
	//	if err != nil {
	//		log.Error().Msgf("Caught an error due closing repository:%+v", err)
	//	}
	//
	//	log.Info().Msg("Everything is closed properly")
	//	cancel()
	//}()
	//if err := srv.Shutdown(ctx); err != nil {
	//	log.Error().Msgf("Server Shutdown Failed:%+v", err)
	//}
	//stop()
	//log.Info().Msg("Server exited properly")

	//log.Info().Msg("Server stopped")
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer func() {
	//	err := repo.Close()
	//	if err != nil {
	//		log.Error().Msgf("Caught an error due closing repository:%+v", err)
	//	}
	//
	//	log.Info().Msg("Everything is closed properly")
	//	cancel()
	//}()
	//if err := srv.Shutdown(ctx); err != nil {
	//	log.Error().Msgf("Server Shutdown Failed:%+v", err)
	//}
	//stop()
	//log.Info().Msg("Server exited properly")
}
