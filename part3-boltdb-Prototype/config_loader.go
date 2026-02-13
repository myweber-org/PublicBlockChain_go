
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	Database string `json:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Features []string       `json:"features"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	var config AppConfig

	if configPath != "" {
		fileData, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := json.Unmarshal(fileData, &config); err != nil {
			return nil, fmt.Errorf("failed to parse config JSON: %w", err)
		}
	}

	overrideFromEnv(&config)

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) {
	overrideStruct(&config.Server)
	overrideStruct(&config.Database)
}

func overrideStruct(target interface{}) {
	// This would use reflection to read struct tags
	// and override values from environment variables
	// Simplified implementation for demonstration
}

func validateConfig(config *AppConfig) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", config.Database.Port)
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[strings.ToLower(config.Server.LogLevel)] {
		return fmt.Errorf("invalid log level: %s", config.Server.LogLevel)
	}

	return nil
}

func (c *AppConfig) String() string {
	// Hide sensitive information
	safeConfig := *c
	safeConfig.Database.Password = "***REDACTED***"
	safeConfig.Database.Username = "***REDACTED***"

	data, _ := json.MarshalIndent(safeConfig, "", "  ")
	return string(data)
}