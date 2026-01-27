package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fiozap/fiozap/internal/api/router"
	"github.com/fiozap/fiozap/internal/config"
	"github.com/fiozap/fiozap/internal/database"
	"github.com/fiozap/fiozap/internal/logger"
	"github.com/fiozap/fiozap/internal/zap"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, cfg.LogFormat)
	log.Info().Msg("Starting FioZap API")

	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer func() { _ = db.Close() }()
	log.Info().Msg("Connected to database")

	manager := zap.NewManager(db.Container, log)

	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	server := &http.Server{
		Addr:    addr,
		Handler: router.New(manager, log, cfg.GlobalAPIToken),
	}

	go func() {
		log.Info().Str("address", addr).Msg("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to gracefully shutdown server")
	}

	log.Info().Msg("Server stopped")
}
