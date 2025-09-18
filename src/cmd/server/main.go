package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"country-iso-matcher/src/internal/config"
	"country-iso-matcher/src/internal/factory"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.Load()

	// Create application factory
	appFactory := factory.NewApplicationFactory(cfg, logger)

	// Build and start server
	server := appFactory.CreateHTTPServer()

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		logger.Info("shutting down server...")
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown failed", "error", err)
		}
		cancel()
	}()

	logger.Info("starting server", "port", cfg.Port)
	if err := server.Start(); err != nil {
		logger.Error("server failed to start", "error", err)
		os.Exit(1)
	}

	<-ctx.Done()
	logger.Info("server stopped")
}
