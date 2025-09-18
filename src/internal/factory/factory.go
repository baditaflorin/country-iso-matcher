package factory

import (
	"log/slog"

	"country-iso-matcher/src/internal/config"
	"country-iso-matcher/src/internal/handler"
	"country-iso-matcher/src/internal/repository/memory"
	"country-iso-matcher/src/internal/server"
	"country-iso-matcher/src/internal/service"
	"country-iso-matcher/src/pkg/normalizer"
)

type ApplicationFactory struct {
	config *config.Config
	logger *slog.Logger
}

func NewApplicationFactory(config *config.Config, logger *slog.Logger) *ApplicationFactory {
	return &ApplicationFactory{
		config: config,
		logger: logger,
	}
}

func (f *ApplicationFactory) CreateHTTPServer() server.Server {
	textNormalizer := normalizer.NewTextNormalizer()
	countryRepo := memory.NewCountryRepository(textNormalizer)
	countryService := service.NewCountryService(countryRepo)
	countryHandler := handler.NewCountryHandler(countryService, f.logger)

	return server.NewHTTPServer(f.config, countryHandler, f.logger)
}
