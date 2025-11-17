package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Load loads configuration from file and environment variables
// Environment variables override file configuration
func Load(configPath string) (*Config, error) {
	// Start with defaults
	cfg := DefaultConfig()

	// Load from file if path provided
	if configPath != "" {
		if err := loadFromFile(cfg, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with environment variables
	loadFromEnv(cfg)

	// Validate configuration
	if err := Validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// loadFromFile loads configuration from a YAML file
func loadFromFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// loadFromEnv overrides configuration with environment variables
func loadFromEnv(cfg *Config) {
	// Server configuration
	if v := os.Getenv("SERVER_PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("PORT"); v != "" { // Backward compatibility
		cfg.Server.Port = v
	}
	if v := os.Getenv("SERVER_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("SERVER_ENVIRONMENT"); v != "" {
		cfg.Server.Environment = v
	}
	if v := os.Getenv("ENV"); v != "" { // Backward compatibility
		cfg.Server.Environment = v
	}
	if v := os.Getenv("SERVER_READ_TIMEOUT"); v != "" {
		if timeout, err := strconv.Atoi(v); err == nil {
			cfg.Server.ReadTimeout = timeout
		}
	}
	if v := os.Getenv("READ_TIMEOUT"); v != "" { // Backward compatibility
		if timeout, err := strconv.Atoi(v); err == nil {
			cfg.Server.ReadTimeout = timeout
		}
	}
	if v := os.Getenv("SERVER_WRITE_TIMEOUT"); v != "" {
		if timeout, err := strconv.Atoi(v); err == nil {
			cfg.Server.WriteTimeout = timeout
		}
	}
	if v := os.Getenv("WRITE_TIMEOUT"); v != "" { // Backward compatibility
		if timeout, err := strconv.Atoi(v); err == nil {
			cfg.Server.WriteTimeout = timeout
		}
	}

	// Database configuration
	if v := os.Getenv("DB_ENABLED"); v != "" {
		cfg.Database.Enabled = v == "true" || v == "1"
	}
	if v := os.Getenv("DB_TYPE"); v != "" {
		cfg.Database.Type = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Database.Port = port
		}
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Database.Database = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.Database.Username = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("DB_SSLMODE"); v != "" {
		cfg.Database.SSLMode = v
	}

	// Database schema configuration
	if v := os.Getenv("DB_COUNTRIES_TABLE"); v != "" {
		cfg.Database.Schema.CountriesTable = v
	}
	if v := os.Getenv("DB_ALIASES_TABLE"); v != "" {
		cfg.Database.Schema.AliasesTable = v
	}
	if v := os.Getenv("DB_CODE_COLUMN"); v != "" {
		cfg.Database.Schema.CodeColumn = v
	}
	if v := os.Getenv("DB_NAME_COLUMN"); v != "" {
		cfg.Database.Schema.NameColumn = v
	}
	if v := os.Getenv("DB_ALIAS_CODE_COLUMN"); v != "" {
		cfg.Database.Schema.AliasCodeColumn = v
	}
	if v := os.Getenv("DB_ALIAS_NAME_COLUMN"); v != "" {
		cfg.Database.Schema.AliasNameColumn = v
	}

	// Data source configuration
	if v := os.Getenv("DATA_SOURCE"); v != "" {
		cfg.Data.Source = v
	}
	if v := os.Getenv("DATA_COUNTRIES_FILE"); v != "" {
		cfg.Data.CountriesFile = v
	}
	if v := os.Getenv("DATA_ALIASES_FILE"); v != "" {
		cfg.Data.AliasesFile = v
	}

	// Logging configuration
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Logging.Level = v
	}
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		cfg.Logging.Format = v
	}

	// GUI configuration
	if v := os.Getenv("GUI_ENABLED"); v != "" {
		cfg.GUI.Enabled = v == "true" || v == "1"
	}
	if v := os.Getenv("GUI_PATH"); v != "" {
		cfg.GUI.Path = v
	}
}

// Save saves the configuration to a YAML file
func Save(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
