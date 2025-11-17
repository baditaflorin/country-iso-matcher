package gui

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"country-iso-matcher/src/internal/config"
)

// ConfigAPI handles configuration API endpoints
type ConfigAPI struct {
	config     *config.Config
	configPath string
	logger     *slog.Logger
}

// NewConfigAPI creates a new configuration API handler
func NewConfigAPI(cfg *config.Config, configPath string, logger *slog.Logger) *ConfigAPI {
	return &ConfigAPI{
		config:     cfg,
		configPath: configPath,
		logger:     logger,
	}
}

// GetConfig returns the current configuration as JSON
func (api *ConfigAPI) GetConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(api.config); err != nil {
		api.logger.Error("Failed to encode config", "error", err)
		http.Error(w, "Failed to encode configuration", http.StatusInternalServerError)
		return
	}
}

// SaveConfig saves the provided configuration
func (api *ConfigAPI) SaveConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newConfig config.Config
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		api.logger.Error("Failed to decode config", "error", err)
		http.Error(w, "Invalid configuration format", http.StatusBadRequest)
		return
	}

	// Validate the new configuration
	if err := config.Validate(&newConfig); err != nil {
		api.logger.Error("Invalid config", "error", err)
		http.Error(w, "Invalid configuration: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Determine config path
	savePath := api.configPath
	if savePath == "" {
		savePath = "config.yaml"
	}

	// Save to file
	if err := config.Save(&newConfig, savePath); err != nil {
		api.logger.Error("Failed to save config", "error", err)
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	api.logger.Info("Configuration saved", "path", savePath)

	// Update current config
	api.config = &newConfig

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Configuration saved successfully",
		"path":    savePath,
	})
}

// ReloadConfig reloads the configuration from file
func (api *ConfigAPI) ReloadConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if api.configPath == "" {
		http.Error(w, "No configuration file specified", http.StatusBadRequest)
		return
	}

	// Check if file exists
	if _, err := os.Stat(api.configPath); os.IsNotExist(err) {
		http.Error(w, "Configuration file not found", http.StatusNotFound)
		return
	}

	// Reload configuration
	newConfig, err := config.Load(api.configPath)
	if err != nil {
		api.logger.Error("Failed to reload config", "error", err)
		http.Error(w, "Failed to reload configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	api.config = newConfig
	api.logger.Info("Configuration reloaded", "path", api.configPath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Configuration reloaded successfully",
	})
}
