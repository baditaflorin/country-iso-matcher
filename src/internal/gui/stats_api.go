package gui

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"

	"country-iso-matcher/src/internal/domain"
	"github.com/prometheus/client_golang/prometheus"
)

// StatsAPI handles statistics API endpoints for the GUI
type StatsAPI struct {
	logger *slog.Logger
}

// NewStatsAPI creates a new stats API handler
func NewStatsAPI(logger *slog.Logger) *StatsAPI {
	return &StatsAPI{
		logger: logger,
	}
}

type StatsResponse struct {
	TotalRequests        float64          `json:"total_requests"`
	SuccessCount         float64          `json:"success_count"`
	NotFoundCount        float64          `json:"not_found_count"`
	ErrorCount           float64          `json:"error_count"`
	ValidationErrorCount float64          `json:"validation_error_count"`
	SuccessRate          float64          `json:"success_rate"`
	FailureRate          float64          `json:"failure_rate"`
	PopularCountries     []PopularCountry `json:"popular_countries,omitempty"`
}

type PopularCountry struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Count float64 `json:"count"`
}

// GetStats returns statistics about country lookups
func (api *StatsAPI) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := StatsResponse{}

	// Get metrics from Prometheus registry
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		api.logger.Error("failed to gather metrics", "error", err)
		api.handleError(w, domain.NewInternalError("Failed to gather metrics"))
		return
	}

	var popularCountries []PopularCountry

	// Parse metrics
	for _, mf := range metricFamilies {
		switch mf.GetName() {
		case "country_lookups_total":
			for _, metric := range mf.GetMetric() {
				value := metric.GetCounter().GetValue()
				for _, label := range metric.GetLabel() {
					if label.GetName() == "result" {
						switch label.GetValue() {
						case "success":
							stats.SuccessCount = value
						case "not_found":
							stats.NotFoundCount = value
						case "error":
							stats.ErrorCount = value
						case "validation_error":
							stats.ValidationErrorCount = value
						}
					}
				}
			}
		case "popular_countries_total":
			for _, metric := range mf.GetMetric() {
				var country PopularCountry
				country.Count = metric.GetCounter().GetValue()
				for _, label := range metric.GetLabel() {
					switch label.GetName() {
					case "country_code":
						country.Code = label.GetValue()
					case "country_name":
						country.Name = label.GetValue()
					}
				}
				if country.Count > 0 {
					popularCountries = append(popularCountries, country)
				}
			}
		}
	}

	// Sort popular countries by count (descending)
	sort.Slice(popularCountries, func(i, j int) bool {
		return popularCountries[i].Count > popularCountries[j].Count
	})

	// Take top 10 most popular countries
	if len(popularCountries) > 10 {
		popularCountries = popularCountries[:10]
	}
	stats.PopularCountries = popularCountries

	// Calculate totals and rates
	stats.TotalRequests = stats.SuccessCount + stats.NotFoundCount + stats.ErrorCount + stats.ValidationErrorCount
	if stats.TotalRequests > 0 {
		stats.SuccessRate = stats.SuccessCount / stats.TotalRequests
		stats.FailureRate = (stats.NotFoundCount + stats.ErrorCount + stats.ValidationErrorCount) / stats.TotalRequests
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		api.logger.Error("failed to encode stats response", "error", err)
	}
}

func (api *StatsAPI) handleError(w http.ResponseWriter, err error) {
	appErr, ok := err.(*domain.AppError)
	if !ok {
		appErr = domain.NewInternalError("Internal server error")
		api.logger.Error("unexpected error", "error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	json.NewEncoder(w).Encode(appErr)
}
