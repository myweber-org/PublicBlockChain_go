package config

import (
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	CacheSize   int
}

func LoadConfig() (*AppConfig, error) {
	config := &AppConfig{
		ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
		DebugMode:   getEnvAsBool("DEBUG_MODE", false),
		DatabaseURL: getEnv("DATABASE_URL", "localhost:5432"),
		CacheSize:   getEnvAsInt("CACHE_SIZE", 100),
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func validateConfig(config *AppConfig) error {
	if config.ServerPort < 1 || config.ServerPort > 65535 {
		return &ConfigError{Field: "ServerPort", Message: "port must be between 1 and 65535"}
	}
	if config.DatabaseURL == "" {
		return &ConfigError{Field: "DatabaseURL", Message: "database URL cannot be empty"}
	}
	if config.CacheSize < 0 {
		return &ConfigError{Field: "CacheSize", Message: "cache size cannot be negative"}
	}
	return nil
}

type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return strings.Join([]string{e.Field, e.Message}, ": ")
}