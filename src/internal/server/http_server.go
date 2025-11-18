package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"country-iso-matcher/src/internal/config"
	"country-iso-matcher/src/internal/gui"
	"country-iso-matcher/src/internal/handler"
	"country-iso-matcher/src/internal/handler/middleware"
	"country-iso-matcher/src/internal/metrics"
	"country-iso-matcher/src/internal/service"
)

type httpServer struct {
	server  *http.Server
	handler handler.CountryHandler
	logger  *slog.Logger
}

func NewHTTPServer(cfg *config.Config, countryHandler handler.CountryHandler, countryService service.CountryService, logger *slog.Logger) Server {
	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/convert", countryHandler.ConvertCountry)
	mux.HandleFunc("/health", countryHandler.Health)
	mux.HandleFunc("/stats", countryHandler.GetStats)
	mux.Handle("/metrics", promhttp.Handler()) // Prometheus metrics endpoint

	// GUI Routes (if enabled)
	if cfg.GUI.Enabled {
		guiHandler := gui.NewHandler(logger)
		configAPI := gui.NewConfigAPI(cfg, "", logger)
		lookupAPI := gui.NewLookupAPI(countryService, logger)
		statsAPI := gui.NewStatsAPI(logger)

		// Serve GUI static files
		guiPath := cfg.GUI.Path
		if !strings.HasSuffix(guiPath, "/") {
			guiPath += "/"
		}

		mux.HandleFunc(guiPath, func(w http.ResponseWriter, r *http.Request) {
			// Remove the GUI path prefix and serve
			r.URL.Path = strings.TrimPrefix(r.URL.Path, strings.TrimSuffix(guiPath, "/"))
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
			guiHandler.ServeGUI(w, r)
		})

		// Config API endpoints
		mux.HandleFunc("/api/config", configAPI.GetConfig)
		mux.HandleFunc("/api/config/save", configAPI.SaveConfig)
		mux.HandleFunc("/api/config/reload", configAPI.ReloadConfig)

		// Lookup and Stats API endpoints for GUI
		mux.HandleFunc("/api/gui/lookup", lookupAPI.LookupCountry)
		mux.HandleFunc("/api/gui/stats", statsAPI.GetStats)

		logger.Info("GUI enabled", "path", cfg.GUI.Path)
	}

	// Root route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("Country ISO Matcher API. Use /api/convert?country=YourCountryName"))
	})

	// Apply middleware (order matters!)
	var httpHandler http.Handler = mux
	httpHandler = middleware.CORS(httpHandler)
	httpHandler = middleware.PrometheusMetrics(httpHandler) // Add Prometheus metrics
	httpHandler = middleware.Logging(logger)(httpHandler)
	httpHandler = middleware.Recovery(logger)(httpHandler)

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	server := &http.Server{
		Addr:         addr,
		Handler:      httpHandler,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Set build info
	metrics.SetBuildInfo("1.0.0", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	// Start system metrics collection (only custom metrics, not conflicting ones)
	metrics.StartSystemMetricsCollection()

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
