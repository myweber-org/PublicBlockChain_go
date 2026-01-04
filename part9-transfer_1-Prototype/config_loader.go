package config

import (
	"os"
	"strconv"
	"strings"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Debug    bool
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.Database.Username = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.Database = getEnv("DB_NAME", "appdb")

	cfg.Server.Port = getEnvAsInt("SERVER_PORT", 8080)
	cfg.Server.ReadTimeout = getEnvAsInt("READ_TIMEOUT", 30)
	cfg.Server.WriteTimeout = getEnvAsInt("WRITE_TIMEOUT", 30)

	cfg.Debug = getEnvAsBool("DEBUG", false)

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
	if valueStr == "" {
		return defaultValue
	}
	return strings.ToLower(valueStr) == "true" || valueStr == "1"
}

func validateConfig(cfg *Config) error {
	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return &ConfigError{Field: "DB_PORT", Message: "port must be between 1 and 65535"}
	}

	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return &ConfigError{Field: "SERVER_PORT", Message: "port must be between 1 and 65535"}
	}

	if cfg.Database.Host == "" {
		return &ConfigError{Field: "DB_HOST", Message: "host cannot be empty"}
	}

	if cfg.Database.Database == "" {
		return &ConfigError{Field: "DB_NAME", Message: "database name cannot be empty"}
	}

	return nil
}

type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return "config error: " + e.Field + " - " + e.Message
}