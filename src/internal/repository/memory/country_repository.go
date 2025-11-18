package memory

import (
	"fmt"

	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/internal/domain"
	"country-iso-matcher/src/pkg/normalizer"
)

type countryRepository struct {
	nameToCode    map[string]string
	codeToCountry map[string]*domain.Country
	normalizer    normalizer.TextNormalizer
}

// NewCountryRepository creates a new in-memory country repository
// It uses a data loader to load country data from various sources (CSV, TSV, memory, database)
func NewCountryRepository(normalizer normalizer.TextNormalizer, loader data.Loader) (*countryRepository, error) {
	repo := &countryRepository{
		nameToCode:    make(map[string]string),
		codeToCountry: make(map[string]*domain.Country),
		normalizer:    normalizer,
	}

	if err := repo.loadCountries(loader); err != nil {
		return nil, fmt.Errorf("failed to load country data: %w", err)
	}

	return repo, nil
}

// FindByName finds a country by its name (supports aliases and fuzzy matching)
func (r *countryRepository) FindByName(name string) (*domain.Country, error) {
	normalized := r.normalizer.Normalize(name)
	code, exists := r.nameToCode[normalized]
	if !exists {
		return nil, domain.NewNotFoundError(name)
	}
	return r.codeToCountry[code], nil
}

// FindByCode finds a country by its ISO code
func (r *countryRepository) FindByCode(code string) (*domain.Country, error) {
	country, exists := r.codeToCountry[code]
	if !exists {
		return nil, domain.NewNotFoundError(code)
	}
	return country, nil
}

// loadCountries loads country data and aliases from the data loader
func (r *countryRepository) loadCountries(loader data.Loader) error {
	// Load countries
	countries, err := loader.LoadCountries()
	if err != nil {
		return fmt.Errorf("failed to load countries: %w", err)
	}

	// Load aliases
	aliases, err := loader.LoadAliases()
	if err != nil {
		return fmt.Errorf("failed to load aliases: %w", err)
	}

	// Build country code map and add all names to lookup
	for i := range countries {
		country := &countries[i]

		// Store by both ISO2 and ISO3 codes
		r.codeToCountry[country.ISO2] = country
		r.codeToCountry[country.ISO3] = country

		// Add all multilingual names to lookup map
		for _, name := range country.Names {
			normalized := r.normalizer.Normalize(name)
			r.nameToCode[normalized] = country.ISO2
		}

		// Add ISO codes themselves as lookup keys
		r.nameToCode[r.normalizer.Normalize(country.ISO2)] = country.ISO2
		r.nameToCode[r.normalizer.Normalize(country.ISO3)] = country.ISO2
	}

	// Build alias lookup map
	for isoCode, aliasNames := range aliases {
		for _, alias := range aliasNames {
			normalized := r.normalizer.Normalize(alias)
			r.nameToCode[normalized] = isoCode
		}
	}

	return nil
}
