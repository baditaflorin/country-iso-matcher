package data

import (
	"country-iso-matcher/src/internal/domain"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// JSONLoader loads country data from individual JSON files
type JSONLoader struct {
	countriesDir string
}

// NewJSONLoader creates a new JSON loader
func NewJSONLoader(countriesDir string) *JSONLoader {
	return &JSONLoader{
		countriesDir: countriesDir,
	}
}

// LoadCountries loads all country JSON files from the directory
func (l *JSONLoader) LoadCountries() ([]domain.Country, error) {
	var countries []domain.Country

	// Read all JSON files from the countries directory
	files, err := filepath.Glob(filepath.Join(l.countriesDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob country files: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no JSON country files found in %s", l.countriesDir)
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", file, err)
		}

		var country domain.Country
		if err := json.Unmarshal(data, &country); err != nil {
			return nil, fmt.Errorf("failed to parse JSON in %s: %w", file, err)
		}

		// Validate required fields
		if country.ISO2 == "" {
			return nil, fmt.Errorf("missing ISO2 code in %s", file)
		}
		if country.ISO3 == "" {
			return nil, fmt.Errorf("missing ISO3 code in %s", file)
		}
		if len(country.Names) == 0 {
			return nil, fmt.Errorf("missing names in %s", file)
		}

		countries = append(countries, country)
	}

	return countries, nil
}

// LoadAliases extracts aliases from country data (for compatibility with existing interface)
func (l *JSONLoader) LoadAliases() (map[string][]string, error) {
	countries, err := l.LoadCountries()
	if err != nil {
		return nil, err
	}

	aliases := make(map[string][]string)
	for _, country := range countries {
		if len(country.Aliases) > 0 {
			aliases[country.ISO2] = country.Aliases
		}
	}

	return aliases, nil
}
