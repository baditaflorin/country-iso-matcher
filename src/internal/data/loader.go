package data

import "country-iso-matcher/src/internal/domain"

// Loader defines the interface for loading country data from various sources
type Loader interface {
	// LoadCountries loads the list of countries with their ISO codes
	LoadCountries() ([]domain.Country, error)

	// LoadAliases loads country name aliases mapped to ISO codes
	LoadAliases() (map[string][]string, error)
}

// CountryData represents the complete country dataset
type CountryData struct {
	Countries []domain.Country
	Aliases   map[string][]string
}

// Load is a helper function to load both countries and aliases
func Load(loader Loader) (*CountryData, error) {
	countries, err := loader.LoadCountries()
	if err != nil {
		return nil, err
	}

	aliases, err := loader.LoadAliases()
	if err != nil {
		return nil, err
	}

	return &CountryData{
		Countries: countries,
		Aliases:   aliases,
	}, nil
}
