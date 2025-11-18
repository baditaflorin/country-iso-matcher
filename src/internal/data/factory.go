package data

import (
	"fmt"

	"country-iso-matcher/src/internal/config"
)

// NewLoader creates the appropriate data loader based on configuration
func NewLoader(cfg *config.DataConfig) (Loader, error) {
	switch cfg.Source {
	case "memory":
		return NewMemoryLoader(), nil

	case "json":
		if cfg.CountriesDir == "" {
			return nil, fmt.Errorf("countries_dir must be specified for JSON source")
		}
		return NewJSONLoader(cfg.CountriesDir), nil

	case "csv":
		if cfg.CountriesFile == "" || cfg.AliasesFile == "" {
			return nil, fmt.Errorf("countries_file and aliases_file must be specified for CSV source")
		}
		return NewCSVLoader(cfg.CountriesFile, cfg.AliasesFile), nil

	case "tsv":
		if cfg.CountriesFile == "" || cfg.AliasesFile == "" {
			return nil, fmt.Errorf("countries_file and aliases_file must be specified for TSV source")
		}
		return NewTSVLoader(cfg.CountriesFile, cfg.AliasesFile), nil

	case "database":
		return nil, fmt.Errorf("database loader not yet implemented - will be available with database repository")

	default:
		return nil, fmt.Errorf("unknown data source: %s (must be json, memory, csv, tsv, or database)", cfg.Source)
	}
}
