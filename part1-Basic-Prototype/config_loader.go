package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    DatabaseURL  string
    MaxConnections int
    DebugMode    bool
    AllowedHosts []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
        DebugMode:    getEnvAsBool("DEBUG_MODE", false),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost", "127.0.0.1"}),
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
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    value, err := strconv.Atoi(strValue)
    if err != nil {
        return defaultValue
    }
    return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    value, err := strconv.ParseBool(strValue)
    if err != nil {
        return defaultValue
    }
    return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    return strings.Split(strValue, ",")
}package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int    `json:"server_port"`
	DebugMode  bool   `json:"debug_mode"`
	Database   string `json:"database_url"`
	CacheTTL   int    `json:"cache_ttl"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	config := &AppConfig{
		ServerPort: 8080,
		DebugMode:  false,
		Database:   "postgres://localhost:5432/appdb",
		CacheTTL:   300,
	}

	if configPath != "" {
		fileData, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := json.Unmarshal(fileData, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	loadFromEnv(config)

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func loadFromEnv(config *AppConfig) {
	if portStr := os.Getenv("APP_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil && port > 0 {
			config.ServerPort = port
		}
	}

	if debugStr := os.Getenv("APP_DEBUG"); debugStr != "" {
		config.DebugMode = strings.ToLower(debugStr) == "true"
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.Database = dbURL
	}

	if ttlStr := os.Getenv("CACHE_TTL"); ttlStr != "" {
		if ttl, err := strconv.Atoi(ttlStr); err == nil && ttl > 0 {
			config.CacheTTL = ttl
		}
	}
}

func validateConfig(config *AppConfig) error {
	if config.ServerPort < 1 || config.ServerPort > 65535 {
		return fmt.Errorf("invalid server port: %d", config.ServerPort)
	}

	if config.Database == "" {
		return fmt.Errorf("database URL cannot be empty")
	}

	if config.CacheTTL < 0 {
		return fmt.Errorf("cache TTL cannot be negative")
	}

	return nil
}