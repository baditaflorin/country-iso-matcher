package gui

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"country-iso-matcher/src/internal/domain"
	"country-iso-matcher/src/internal/service"
)

// LookupAPI handles country lookup API endpoints for the GUI
type LookupAPI struct {
	service service.CountryService
	logger  *slog.Logger
}

// NewLookupAPI creates a new lookup API handler
func NewLookupAPI(service service.CountryService, logger *slog.Logger) *LookupAPI {
	return &LookupAPI{
		service: service,
		logger:  logger,
	}
}

// LookupCountry performs a country lookup
func (api *LookupAPI) LookupCountry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	countryName := r.URL.Query().Get("country")
	if countryName == "" {
		api.handleError(w, domain.NewValidationError("Country name is required", ""), "")
		return
	}

	result, err := api.service.LookupCountry(countryName)
	if err != nil {
		api.handleError(w, err, countryName)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		api.logger.Error("failed to encode response", "error", err)
	}
}

func (api *LookupAPI) handleError(w http.ResponseWriter, err error, query string) {
	appErr, ok := err.(*domain.AppError)
	if !ok {
		appErr = domain.NewInternalError("Internal server error")
		api.logger.Error("unexpected error", "error", err, "query", query)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	json.NewEncoder(w).Encode(appErr)
}
