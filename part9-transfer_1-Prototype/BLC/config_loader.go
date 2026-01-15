
package config

import (
	"encoding/json"
	"os"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Features []string       `json:"features"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	config := &AppConfig{}
	
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	overrideWithEnvVars(config)
	
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func overrideWithEnvVars(config *AppConfig) {
	overrideStruct(&config.Database)
	overrideStruct(&config.Server)
}

func overrideStruct(s interface{}) {
	// Implementation would use reflection to check struct tags
	// and override values from environment variables
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return NewValidationError("database host cannot be empty")
	}
	
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return NewValidationError("database port must be between 1 and 65535")
	}
	
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return NewValidationError("server port must be between 1 and 65535")
	}
	
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	
	if !validLogLevels[strings.ToLower(config.Server.LogLevel)] {
		return NewValidationError("invalid log level specified")
	}
	
	return nil
}

type ValidationError struct {
	Message string
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{Message: msg}
}

func (e *ValidationError) Error() string {
	return "config validation error: " + e.Message
}package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	ServerPort string `json:"server_port"`
	DatabaseURL string `json:"database_url"`
	LogLevel string `json:"log_level"`
	CacheEnabled bool `json:"cache_enabled"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	if port := os.Getenv("APP_PORT"); port != "" {
		config.ServerPort = port
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.DatabaseURL = dbURL
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	return &config, nil
}package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DBHost     string
    DBPort     int
    DebugMode  bool
    APIKeys    []string
}

func Load() (*Config, error) {
    cfg := &Config{}
    var err error

    cfg.ServerPort, err = getIntEnv("SERVER_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }

    cfg.DBHost = getStringEnv("DB_HOST", "localhost")
    
    cfg.DBPort, err = getIntEnv("DB_PORT", 5432)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }

    cfg.DebugMode, err = getBoolEnv("DEBUG_MODE", false)
    if err != nil {
        return nil, fmt.Errorf("invalid DEBUG_MODE: %w", err)
    }

    cfg.APIKeys = getStringSliceEnv("API_KEYS", []string{})

    if err := validateConfig(cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return cfg, nil
}

func getStringEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
    if value := os.Getenv(key); value != "" {
        intValue, err := strconv.Atoi(value)
        if err != nil {
            return 0, fmt.Errorf("cannot parse %s as integer: %w", key, err)
        }
        return intValue, nil
    }
    return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) (bool, error) {
    if value := os.Getenv(key); value != "" {
        boolValue, err := strconv.ParseBool(value)
        if err != nil {
            return false, fmt.Errorf("cannot parse %s as boolean: %w", key, err)
        }
        return boolValue, nil
    }
    return defaultValue, nil
}

func getStringSliceEnv(key string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        return strings.Split(value, ",")
    }
    return defaultValue
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    
    if cfg.DBPort < 1 || cfg.DBPort > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }
    
    if cfg.DBHost == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    
    return nil
}