package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"country-iso-matcher/src/internal/domain"
	"country-iso-matcher/src/internal/service"
)

type countryHandler struct {
	service service.CountryService
	logger  *slog.Logger
}

func NewCountryHandler(service service.CountryService, logger *slog.Logger) CountryHandler {
	return &countryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *countryHandler) ConvertCountry(w http.ResponseWriter, r *http.Request) {
	countryName := r.URL.Query().Get("country")

	result, err := h.service.LookupCountry(countryName)
	if err != nil {
		h.handleError(w, err, countryName)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *countryHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "country-iso-matcher",
	})
}

func (h *countryHandler) handleError(w http.ResponseWriter, err error, query string) {
	appErr, ok := err.(*domain.AppError)
	if !ok {
		appErr = domain.NewInternalError("Internal server error")
		h.logger.Error("unexpected error", "error", err, "query", query)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	json.NewEncoder(w).Encode(appErr)
}
