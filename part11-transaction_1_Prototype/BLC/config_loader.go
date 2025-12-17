package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	AllowedIPs []string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.ServerPort = getEnvAsInt("SERVER_PORT", 8080)
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort = getEnvAsInt("DB_PORT", 5432)
	cfg.DebugMode = getEnvAsBool("DEBUG_MODE", false)
	cfg.AllowedIPs = getEnvAsSlice("ALLOWED_IPS", []string{"127.0.0.1"}, ",")

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
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

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, sep)
}

func validateConfig(cfg *Config) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return &ConfigError{Field: "SERVER_PORT", Value: cfg.ServerPort}
	}
	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return &ConfigError{Field: "DB_PORT", Value: cfg.DBPort}
	}
	return nil
}

type ConfigError struct {
	Field string
	Value interface{}
}

func (e *ConfigError) Error() string {
	return "invalid configuration value for field: " + e.Field
}