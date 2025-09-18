package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"country-iso-matcher/src/internal/config"
	"country-iso-matcher/src/internal/handler"
	"country-iso-matcher/src/internal/handler/middleware"
)

type httpServer struct {
	server  *http.Server
	handler handler.CountryHandler
	logger  *slog.Logger
}

func NewHTTPServer(cfg *config.Config, countryHandler handler.CountryHandler, logger *slog.Logger) Server {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/api/convert", countryHandler.ConvertCountry)
	mux.HandleFunc("/health", countryHandler.Health)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Country ISO Matcher API. Use /api/convert?country=YourCountryName"))
	})

	// Apply middleware
	var handler http.Handler = mux
	handler = middleware.CORS(handler)
	handler = middleware.Logging(logger)(handler)
	handler = middleware.Recovery(logger)(handler)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}

	return &httpServer{
		server:  server,
		handler: countryHandler,
		logger:  logger,
	}
}

func (s *httpServer) Start() error {
	s.logger.Info("starting HTTP server", "addr", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed to start: %w", err)
	}
	return nil
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down HTTP server")
	return s.server.Shutdown(ctx)
}
