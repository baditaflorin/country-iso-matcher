package service

import (
	"strings"

	"country-iso-matcher/src/internal/domain"
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
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, domain.NewValidationError("Country query parameter is required", query)
	}

	country, err := s.repository.FindByName(query)
	if err != nil {
		return nil, err
	}

	return domain.NewCountryResponse(query, country.Name, country.Code), nil
}
