package factory

import (
	"fmt"
	"log/slog"

	"country-iso-matcher/src/internal/config"
	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/internal/handler"
	"country-iso-matcher/src/internal/repository/memory"
	"country-iso-matcher/src/internal/server"
	"country-iso-matcher/src/internal/service"
	"country-iso-matcher/src/pkg/normalizer"
)

// ApplicationFactory creates and wires up application dependencies
type ApplicationFactory struct {
	config *config.Config
	logger *slog.Logger
}

// NewApplicationFactory creates a new application factory
func NewApplicationFactory(cfg *config.Config, logger *slog.Logger) (*ApplicationFactory, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	return &ApplicationFactory{
		config: cfg,
		logger: logger,
	}, nil
}

// CreateHTTPServer creates and configures the HTTP server
func (f *ApplicationFactory) CreateHTTPServer() (server.Server, error) {
	// Create data loader based on configuration
	loader, err := data.NewLoader(&f.config.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to create data loader: %w", err)
	}

	// Create text normalizer
	textNormalizer := normalizer.NewTextNormalizer()

	// Create country repository
	countryRepo, err := memory.NewCountryRepository(textNormalizer, loader)
	if err != nil {
		return nil, fmt.Errorf("failed to create country repository: %w", err)
	}

	// Create country service
	countryService := service.NewCountryService(countryRepo)

	// Create country handler
	countryHandler := handler.NewCountryHandler(countryService, f.logger)

	// Create and return HTTP server
	return server.NewHTTPServer(f.config, countryHandler, f.logger), nil
}
