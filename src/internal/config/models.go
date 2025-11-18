package config

// Config represents the complete application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server" json:"server"`
	Database DatabaseConfig `yaml:"database" json:"database"`
	Data     DataConfig     `yaml:"data" json:"data"`
	Logging  LoggingConfig  `yaml:"logging" json:"logging"`
	GUI      GUIConfig      `yaml:"gui" json:"gui"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Port         string `yaml:"port" json:"port"`
	Host         string `yaml:"host" json:"host"`
	Environment  string `yaml:"environment" json:"environment"`
	ReadTimeout  int    `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout" json:"write_timeout"`
}

// DatabaseConfig contains database connection configuration
type DatabaseConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Type     string `yaml:"type" json:"type"` // postgres, mysql, sqlite
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Database string `yaml:"database" json:"database"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	SSLMode  string `yaml:"ssl_mode" json:"ssl_mode"`

	// Configurable table and column names
	Schema SchemaConfig `yaml:"schema" json:"schema"`
}

// SchemaConfig defines database table and column names
type SchemaConfig struct {
	CountriesTable   string `yaml:"countries_table" json:"countries_table"`
	AliasesTable     string `yaml:"aliases_table" json:"aliases_table"`
	CodeColumn       string `yaml:"code_column" json:"code_column"`
	NameColumn       string `yaml:"name_column" json:"name_column"`
	AliasCodeColumn  string `yaml:"alias_code_column" json:"alias_code_column"`
	AliasNameColumn  string `yaml:"alias_name_column" json:"alias_name_column"`
}

// DataConfig specifies the data source configuration
type DataConfig struct {
	Source        string `yaml:"source" json:"source"` // json, memory, csv, tsv, database
	CountriesDir  string `yaml:"countries_dir" json:"countries_dir"` // for JSON source
	CountriesFile string `yaml:"countries_file" json:"countries_file"`
	AliasesFile   string `yaml:"aliases_file" json:"aliases_file"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level" json:"level"`   // debug, info, warn, error
	Format string `yaml:"format" json:"format"` // json, text
}

// GUIConfig contains GUI-related configuration
type GUIConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Path    string `yaml:"path" json:"path"` // URL path for GUI
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         "3030",
			Host:         "0.0.0.0",
			Environment:  "development",
			ReadTimeout:  10,
			WriteTimeout: 10,
		},
		Database: DatabaseConfig{
			Enabled:  false,
			Type:     "postgres",
			Host:     "localhost",
			Port:     5432,
			Database: "countries",
			Username: "postgres",
			Password: "",
			SSLMode:  "disable",
			Schema: SchemaConfig{
				CountriesTable:  "countries",
				AliasesTable:    "country_aliases",
				CodeColumn:      "code",
				NameColumn:      "name",
				AliasCodeColumn: "country_code",
				AliasNameColumn: "alias",
			},
		},
		Data: DataConfig{
			Source:        "json",
			CountriesDir:  "data/countries",
			CountriesFile: "data/countries.csv",
			AliasesFile:   "data/aliases.csv",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		GUI: GUIConfig{
			Enabled: true,
			Path:    "/admin",
		},
	}
}
