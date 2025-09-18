package service

import "country-iso-matcher/src/internal/domain"

type CountryService interface {
	LookupCountry(query string) (*domain.CountryResponse, error)
}
