package service

import (
	"strings"
	"time"

	"country-iso-matcher/src/internal/domain"
	"country-iso-matcher/src/internal/metrics"
	"country-iso-matcher/src/internal/repository"
)

type countryService struct {
	repository repository.CountryRepository
}

func NewCountryService(repo repository.CountryRepository) CountryService {
	return &countryService{
		repository: repo,
	}
}

func (s *countryService) LookupCountry(query string) (*domain.CountryResponse, error) {
	start := time.Now()
	var result string

	defer func() {
		duration := time.Since(start).Seconds()
		metrics.CountryLookupDuration.WithLabelValues(result).Observe(duration)
	}()

	query = strings.TrimSpace(query)
	if query == "" {
		result = "validation_error"
		metrics.CountryLookupsTotal.WithLabelValues("validation_error").Inc()
		return nil, domain.NewValidationError("Country query parameter is required", query)
	}

	country, err := s.repository.FindByName(query)
	if err != nil {
		// Check if it's a not found error or other error
		if appErr, ok := err.(*domain.AppError); ok {
			if appErr.Code == 404 {
				result = "not_found"
				metrics.CountryLookupsTotal.WithLabelValues("not_found").Inc()
			} else {
				result = "error"
				metrics.CountryLookupsTotal.WithLabelValues("error").Inc()
			}
		} else {
			result = "error"
			metrics.CountryLookupsTotal.WithLabelValues("error").Inc()
		}
		return nil, err
	}

	result = "success"
	metrics.CountryLookupsTotal.WithLabelValues("success").Inc()

	// Track popular countries for successful lookups
	metrics.PopularCountries.WithLabelValues(country.ISO2, country.GetOfficialName()).Inc()

	return domain.NewCountryResponse(query, country), nil
}
