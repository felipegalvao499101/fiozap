package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fiozap/internal/api/router"
	"fiozap/internal/config"
	"fiozap/internal/database"
	"fiozap/internal/logger"
	"fiozap/internal/providers/wameow"
	"fiozap/internal/repository"

	_ "fiozap/docs"
)

// @title           FioZap API
// @version         1.0
// @description     API REST multi-session para WhatsApp usando whatsmeow. Baseada na WuzAPI com campos JSON em PascalCase.
// @termsOfService  http://swagger.io/terms/

// @contact.name   FioZap Support
// @contact.url    https://github.com/fiozap/fiozap

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Token global (GLOBAL_API_TOKEN) ou token da sessao

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

	repos := repository.New(db.DB)
	provider := wameow.New(db.Container, repos.Session, log)

	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	server := &http.Server{
		Addr:    addr,
		Handler: router.New(provider, log, cfg.GlobalAPIToken),
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
