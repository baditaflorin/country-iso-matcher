package memory

import (
	"fmt"
	"strings"

	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/internal/domain"
	"country-iso-matcher/src/pkg/normalizer"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
)

type countryRepository struct {
	nameToCode       map[string]string
	codeToCountry    map[string]*domain.Country
	exactNameToCode  map[string]string
	accentNameToCode map[string]string
	normalizer       normalizer.TextNormalizer
}

// NewCountryRepository creates a new in-memory country repository
// It uses a data loader to load country data from various sources (CSV, TSV, memory, database)
func NewCountryRepository(normalizer normalizer.TextNormalizer, loader data.Loader) (*countryRepository, error) {
	repo := &countryRepository{
		nameToCode:       make(map[string]string),
		codeToCountry:    make(map[string]*domain.Country),
		exactNameToCode:  make(map[string]string),
		accentNameToCode: make(map[string]string),
		normalizer:       normalizer,
	}

	if err := repo.loadCountries(loader); err != nil {
		return nil, fmt.Errorf("failed to load country data: %w", err)
	}

	return repo, nil
}

// FindByName finds a country by its name (supports aliases and fuzzy matching)
func (r *countryRepository) FindByName(name string) (*domain.Country, error) {
	trimmed := strings.TrimSpace(name)
	if code, exists := r.exactNameToCode[trimmed]; exists {
		return r.codeToCountry[code], nil
	}
	if containsNonASCII(trimmed) {
		if code, exists := r.accentNameToCode[accentPreservingKey(trimmed)]; exists {
			return r.codeToCountry[code], nil
		}
	}

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
	canonicalNames := make(map[string]string)
	exactAmbiguous := make(map[string]bool)
	accentAmbiguous := make(map[string]bool)
	for i := range countries {
		country := &countries[i]

		// Store by both ISO2 and ISO3 codes
		r.codeToCountry[country.ISO2] = country
		r.codeToCountry[country.ISO3] = country

		// Add all multilingual names to lookup map
		for _, name := range country.Names {
			normalized := r.normalizer.Normalize(name)
			r.nameToCode[normalized] = country.ISO2
			canonicalNames[normalized] = country.ISO2
			addUniqueName(r.exactNameToCode, exactAmbiguous, strings.TrimSpace(name), country.ISO2)
			addUniqueName(r.accentNameToCode, accentAmbiguous, accentPreservingKey(name), country.ISO2)
		}

		// Add ISO codes themselves as lookup keys
		r.nameToCode[r.normalizer.Normalize(country.ISO2)] = country.ISO2
		r.nameToCode[r.normalizer.Normalize(country.ISO3)] = country.ISO2
	}

	// Build alias lookup map
	for isoCode, aliasNames := range aliases {
		for _, alias := range aliasNames {
			normalized := r.normalizer.Normalize(alias)
			if canonicalCode, exists := canonicalNames[normalized]; !exists || canonicalCode == isoCode {
				addUniqueName(r.exactNameToCode, exactAmbiguous, strings.TrimSpace(alias), isoCode)
				addUniqueName(r.accentNameToCode, accentAmbiguous, accentPreservingKey(alias), isoCode)
			}
			// An alias must not shadow a canonical country name or ISO code.
			if _, exists := r.nameToCode[normalized]; exists {
				continue
			}
			r.nameToCode[normalized] = isoCode
		}
	}

	return nil
}

func addUniqueName(names map[string]string, ambiguous map[string]bool, name, isoCode string) {
	if name == "" || ambiguous[name] {
		return
	}
	if existing, exists := names[name]; exists && existing != isoCode {
		delete(names, name)
		ambiguous[name] = true
		return
	}
	names[name] = isoCode
}

func accentPreservingKey(name string) string {
	folded, _, err := transform.String(width.Fold, strings.TrimSpace(name))
	if err != nil {
		return strings.ToLower(strings.TrimSpace(name))
	}
	return norm.NFC.String(strings.ToLower(folded))
}

func containsNonASCII(value string) bool {
	for _, r := range value {
		if r > 0x7f {
			return true
		}
	}
	return false
}
