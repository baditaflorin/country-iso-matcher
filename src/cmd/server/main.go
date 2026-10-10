package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"country-iso-matcher/src/internal/config"
	"country-iso-matcher/src/internal/factory"
	"country-iso-matcher/src/internal/version"

	"github.com/baditaflorin/go-common/profiling"
	"github.com/baditaflorin/go-common/telemetry"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "", "Path to configuration file (YAML)")
	flag.Parse()

	// Load configuration
	// If no config path specified, check CONFIG_FILE env var
	if *configPath == "" {
		*configPath = os.Getenv("CONFIG_FILE")
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Setup structured logging based on configuration
	logger := setupLogger(cfg)
	slog.SetDefault(logger)
	telemetryConfig := telemetry.Init("country-iso-matcher", version.Version)
	stopProfiling, err := profiling.StartFromEnv("country-iso-matcher")
	if err != nil {
		logger.Error("Pyroscope profiling startup failed", "error", err)
	} else {
		defer stopProfiling()
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), telemetry.ExporterTimeout)
		defer cancel()
		if err := telemetryConfig.Shutdown(shutdownCtx); err != nil {
			logger.Error("OpenTelemetry shutdown failed", "error", err)
		}
	}()

	logger.Info("Starting Country ISO Matcher",
		"environment", cfg.Server.Environment,
		"port", cfg.Server.Port,
		"data_source", cfg.Data.Source,
	)

	// Create application factory
	appFactory, err := factory.NewApplicationFactory(cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}

	// Build and start server
	server, err := appFactory.CreateHTTPServer()
	if err != nil {
		logger.Error("Failed to create HTTP server", "error", err)
		os.Exit(1)
	}

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		logger.Info("shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("server shutdown failed", "error", err)
		}
		cancel()
	}()

	logger.Info("server listening", "address", cfg.Server.Host+":"+cfg.Server.Port)
	if err := server.Start(); err != nil {
		logger.Error("server failed to start", "error", err)
		return
	}

	<-ctx.Done()
	logger.Info("server stopped gracefully")
}

// setupLogger creates a logger based on configuration
func setupLogger(cfg *config.Config) *slog.Logger {
	var level slog.Level
	switch cfg.Logging.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if cfg.Logging.Format == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
