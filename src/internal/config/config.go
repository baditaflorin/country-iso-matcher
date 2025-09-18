package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	Environment  string
	ReadTimeout  int
	WriteTimeout int
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "3030"),
		Environment:  getEnv("ENV", "development"),
		ReadTimeout:  getEnvInt("READ_TIMEOUT", 10),
		WriteTimeout: getEnvInt("WRITE_TIMEOUT", 10),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
