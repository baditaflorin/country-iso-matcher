package repository

import "country-iso-matcher/src/internal/domain"

type CountryRepository interface {
	FindByName(name string) (*domain.Country, error)
	FindByCode(code string) (*domain.Country, error)
}
