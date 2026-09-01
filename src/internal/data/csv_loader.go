package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"country-iso-matcher/src/internal/domain"
)

// CSVLoader loads country data from CSV files
type CSVLoader struct {
	countriesFile string
	aliasesFile   string
}

// NewCSVLoader creates a new CSV loader
func NewCSVLoader(countriesFile, aliasesFile string) *CSVLoader {
	return &CSVLoader{
		countriesFile: countriesFile,
		aliasesFile:   aliasesFile,
	}
}

// LoadCountries loads countries from a CSV file
// Expected format: code,iso3,name. The legacy code,name form remains
// supported for callers that have not yet added ISO-3 data.
// Example: US,USA,United States
func (l *CSVLoader) LoadCountries() ([]domain.Country, error) {
	file, err := os.Open(l.countriesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open countries file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	// Allow legacy two-column rows alongside the current three-column form.
	reader.FieldsPerRecord = -1

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read countries CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("countries file is empty")
	}

	countries := make([]domain.Country, 0, len(records))

	// Check if first row is a header
	startIdx := 0
	if isHeader(records[0]) {
		startIdx = 1
	}

	for i := startIdx; i < len(records); i++ {
		record := records[i]
		if len(record) < 2 {
			continue // Skip invalid rows
		}

		code := strings.TrimSpace(record[0])
		iso3 := code
		nameIndex := 1
		if len(record) >= 3 && len(strings.TrimSpace(record[1])) == 3 {
			iso3 = strings.ToUpper(strings.TrimSpace(record[1]))
			nameIndex = 2
		}
		name := strings.TrimSpace(record[nameIndex])

		if code == "" || name == "" {
			continue // Skip empty entries
		}

		countries = append(countries, domain.Country{
			ISO2: code,
			ISO3: iso3, // Legacy two-column data falls back to ISO2.
			Names: map[string]string{
				"en": name,
			},
			Aliases: []string{},
		})
	}

	if len(countries) == 0 {
		return nil, fmt.Errorf("no valid countries found in file")
	}

	return countries, nil
}

// LoadAliases loads country aliases from a CSV file
// Expected format: code,alias1,alias2,alias3,...
// Example: US,usa,america,united states,estados unidos
func (l *CSVLoader) LoadAliases() (map[string][]string, error) {
	file, err := os.Open(l.aliasesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open aliases file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read aliases CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("aliases file is empty")
	}

	aliases := make(map[string][]string)

	// Check if first row is a header
	startIdx := 0
	if isHeader(records[0]) {
		startIdx = 1
	}

	for i := startIdx; i < len(records); i++ {
		record := records[i]
		if len(record) < 2 {
			continue // Skip invalid rows
		}

		code := strings.TrimSpace(record[0])
		if code == "" {
			continue
		}

		// Collect all aliases for this country
		countryAliases := make([]string, 0, len(record)-1)
		for j := 1; j < len(record); j++ {
			alias := strings.TrimSpace(record[j])
			if alias != "" {
				countryAliases = append(countryAliases, alias)
			}
		}

		if len(countryAliases) > 0 {
			aliases[code] = countryAliases
		}
	}

	return aliases, nil
}

// isHeader checks if a record looks like a header row
func isHeader(record []string) bool {
	if len(record) == 0 {
		return false
	}

	// Common header patterns
	firstField := strings.ToLower(strings.TrimSpace(record[0]))
	headers := map[string]bool{
		"code":         true,
		"iso":          true,
		"iso_code":     true,
		"country_code": true,
	}

	return headers[firstField]
}
