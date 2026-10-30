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
	return s.lookupCountry(query, false)
}

// LookupCountryWithUserAssigned enables non-ISO codes only for callers that
// explicitly request them. The default LookupCountry remains ISO-only.
func (s *countryService) LookupCountryWithUserAssigned(query string) (*domain.CountryResponse, error) {
	return s.lookupCountry(query, true)
}

func (s *countryService) lookupCountry(query string, includeUserAssigned bool) (*domain.CountryResponse, error) {
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
	// Prefer the configured country catalog, and use the user-assigned
	// convention only as an explicit fallback for missing entries.
	if appErr, notFound := err.(*domain.AppError); includeUserAssigned && notFound && appErr.Code == 404 {
		switch strings.ToLower(query) {
		case "xk", "xkx", "kosovo", "republic of kosovo":
			result = "success"
			metrics.CountryLookupsTotal.WithLabelValues(result).Inc()
			metrics.PopularCountries.WithLabelValues("XK", "Kosovo").Inc()
			return &domain.CountryResponse{Query: query, OfficialName: "Kosovo", ISO2Code: "XK", ISO3Code: "XKX", UserAssigned: true}, nil
		}
	}
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
