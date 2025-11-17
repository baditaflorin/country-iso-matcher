package config

import (
	"fmt"
	"os"
	"strings"
)

// Validate validates the configuration
func Validate(cfg *Config) error {
	// Validate server configuration
	if err := validateServer(&cfg.Server); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	// Validate database configuration if enabled
	if cfg.Database.Enabled {
		if err := validateDatabase(&cfg.Database); err != nil {
			return fmt.Errorf("database config: %w", err)
		}
	}

	// Validate data source configuration
	if err := validateDataSource(&cfg.Data, cfg.Database.Enabled); err != nil {
		return fmt.Errorf("data source config: %w", err)
	}

	// Validate logging configuration
	if err := validateLogging(&cfg.Logging); err != nil {
		return fmt.Errorf("logging config: %w", err)
	}

	return nil
}

func validateServer(cfg *ServerConfig) error {
	if cfg.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	if cfg.ReadTimeout <= 0 {
		return fmt.Errorf("read_timeout must be positive")
	}

	if cfg.WriteTimeout <= 0 {
		return fmt.Errorf("write_timeout must be positive")
	}

	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
		"test":        true,
	}

	if !validEnvs[cfg.Environment] {
		return fmt.Errorf("invalid environment: %s (must be development, staging, production, or test)", cfg.Environment)
	}

	return nil
}

func validateDatabase(cfg *DatabaseConfig) error {
	validTypes := map[string]bool{
		"postgres": true,
		"mysql":    true,
		"sqlite":   true,
	}

	if !validTypes[cfg.Type] {
		return fmt.Errorf("invalid database type: %s (must be postgres, mysql, or sqlite)", cfg.Type)
	}

	if cfg.Type != "sqlite" {
		if cfg.Host == "" {
			return fmt.Errorf("host cannot be empty for %s", cfg.Type)
		}

		if cfg.Port <= 0 || cfg.Port > 65535 {
			return fmt.Errorf("invalid port: %d", cfg.Port)
		}

		if cfg.Database == "" {
			return fmt.Errorf("database name cannot be empty")
		}

		if cfg.Username == "" {
			return fmt.Errorf("username cannot be empty")
		}
	}

	// Validate schema configuration
	if cfg.Schema.CountriesTable == "" {
		return fmt.Errorf("countries_table cannot be empty")
	}

	if cfg.Schema.AliasesTable == "" {
		return fmt.Errorf("aliases_table cannot be empty")
	}

	if cfg.Schema.CodeColumn == "" {
		return fmt.Errorf("code_column cannot be empty")
	}

	if cfg.Schema.NameColumn == "" {
		return fmt.Errorf("name_column cannot be empty")
	}

	return nil
}

func validateDataSource(cfg *DataConfig, dbEnabled bool) error {
	validSources := map[string]bool{
		"memory":   true,
		"csv":      true,
		"tsv":      true,
		"database": true,
	}

	if !validSources[cfg.Source] {
		return fmt.Errorf("invalid data source: %s (must be memory, csv, tsv, or database)", cfg.Source)
	}

	// If source is database, database must be enabled
	if cfg.Source == "database" && !dbEnabled {
		return fmt.Errorf("data source is 'database' but database is not enabled")
	}

	// Validate file paths for csv/tsv sources
	if cfg.Source == "csv" || cfg.Source == "tsv" {
		if cfg.CountriesFile == "" {
			return fmt.Errorf("countries_file cannot be empty for %s source", cfg.Source)
		}

		if cfg.AliasesFile == "" {
			return fmt.Errorf("aliases_file cannot be empty for %s source", cfg.Source)
		}

		// Check if files exist (warning, not error)
		if _, err := os.Stat(cfg.CountriesFile); os.IsNotExist(err) {
			// This is just a warning - file might be created later
			// But we'll validate in the data loader
		}

		if _, err := os.Stat(cfg.AliasesFile); os.IsNotExist(err) {
			// This is just a warning
		}
	}

	return nil
}

func validateLogging(cfg *LoggingConfig) error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	level := strings.ToLower(cfg.Level)
	if !validLevels[level] {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", cfg.Level)
	}
	cfg.Level = level // Normalize to lowercase

	validFormats := map[string]bool{
		"json": true,
		"text": true,
	}

	format := strings.ToLower(cfg.Format)
	if !validFormats[format] {
		return fmt.Errorf("invalid log format: %s (must be json or text)", cfg.Format)
	}
	cfg.Format = format // Normalize to lowercase

	return nil
}
