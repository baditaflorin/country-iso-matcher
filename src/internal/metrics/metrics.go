package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// Business logic metrics - Enhanced for success/failure tracking
	CountryLookupsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "country_lookups_total",
			Help: "Total number of country lookups by result type",
		},
		[]string{"result"}, // "success", "not_found", "validation_error", "error"
	)

	CountryLookupDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "country_lookup_duration_seconds",
			Help:    "Country lookup duration in seconds by result type",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0},
		},
		[]string{"result"}, // Track duration by success/failure
	)

	// Success rate metrics - easier to query
	CountryLookupSuccessRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "country_lookup_success_rate",
			Help: "Success rate of country lookups (0-1)",
		},
		[]string{"time_window"}, // "1m", "5m", "15m" etc
	)

	// Popular countries metrics
	PopularCountries = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "popular_countries_total",
			Help: "Count of successful lookups by country code",
		},
		[]string{"country_code", "country_name"},
	)

	// System metrics
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	MemoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
	)

	// Application info
	BuildInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "build_info",
			Help: "Build information",
		},
		[]string{"version", "goversion", "goos", "goarch"},
	)
)

// SetBuildInfo sets the build information metrics
func SetBuildInfo(version, goVersion, goos, goarch string) {
	BuildInfo.WithLabelValues(version, goVersion, goos, goarch).Set(1)
}
