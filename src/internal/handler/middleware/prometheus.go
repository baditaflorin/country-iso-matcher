package middleware

import (
	"net/http"
	"strconv"
	"time"

	"country-iso-matcher/src/internal/metrics"
)

type prometheusResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

func (prw *prometheusResponseWriter) WriteHeader(code int) {
	prw.statusCode = code
	prw.ResponseWriter.WriteHeader(code)
}

func (prw *prometheusResponseWriter) Write(b []byte) (int, error) {
	n, err := prw.ResponseWriter.Write(b)
	prw.written += n
	return n, err
}

// PrometheusMetrics middleware collects HTTP metrics for Prometheus
func PrometheusMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		prw := &prometheusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Increment active connections
		metrics.ActiveConnections.Inc()
		defer metrics.ActiveConnections.Dec()

		next.ServeHTTP(prw, r)

		duration := time.Since(start)
		statusCode := strconv.Itoa(prw.statusCode)
		endpoint := getEndpointLabel(r.URL.Path)

		// Record metrics
		metrics.HTTPRequestsTotal.WithLabelValues(
			r.Method,
			endpoint,
			statusCode,
		).Inc()

		metrics.HTTPRequestDuration.WithLabelValues(
			r.Method,
			endpoint,
			statusCode,
		).Observe(duration.Seconds())
	})
}

// getEndpointLabel normalizes endpoint names for metrics
func getEndpointLabel(path string) string {
	switch path {
	case "/api/convert":
		return "convert"
	case "/health":
		return "health"
	case "/metrics":
		return "metrics"
	case "/":
		return "root"
	default:
		return "other"
	}
}
