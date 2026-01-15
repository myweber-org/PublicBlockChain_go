package config

import (
	"errors"
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
	SSLMode  string
}

type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
	DebugMode    bool
}

type AppConfig struct {
	DB     DatabaseConfig
	Server ServerConfig
	APIKey string
}

func LoadConfig() (*AppConfig, error) {
	config := &AppConfig{}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, errors.New("invalid DB_PORT value")
	}

	config.DB = DatabaseConfig{
		Host:     dbHost,
		Port:     dbPort,
		Username: getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASS", ""),
		Database: getEnv("DB_NAME", "appdb"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}

	readTimeout, err := strconv.Atoi(getEnv("READ_TIMEOUT", "30"))
	if err != nil {
		return nil, errors.New("invalid READ_TIMEOUT value")
	}

	writeTimeout, err := strconv.Atoi(getEnv("WRITE_TIMEOUT", "30"))
	if err != nil {
		return nil, errors.New("invalid WRITE_TIMEOUT value")
	}

	debugMode, err := strconv.ParseBool(getEnv("DEBUG_MODE", "false"))
	if err != nil {
		return nil, errors.New("invalid DEBUG_MODE value")
	}

	config.Server = ServerConfig{
		Port:         serverPort,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		DebugMode:    debugMode,
	}

	apiKey := getEnv("API_KEY", "")
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("API_KEY is required")
	}
	config.APIKey = apiKey

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}package config

import (
    "os"
    "strconv"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    DebugMode   bool
    MaxWorkers  int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort: getEnvAsInt("SERVER_PORT", 8080),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        DebugMode:   getEnvAsBool("DEBUG_MODE", false),
        MaxWorkers:  getEnvAsInt("MAX_WORKERS", 10),
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
}package config

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

func Load() (*Config, error) {
    cfg := &Config{
        DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
        DebugMode:    getEnvAsBool("DEBUG_MODE", false),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost"}, ","),
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

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    return strings.Split(strValue, sep)
}