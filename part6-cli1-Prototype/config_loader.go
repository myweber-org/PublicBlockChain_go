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
	FeatureFlags map[string]bool
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		FeatureFlags: make(map[string]bool),
	}

	port, err := getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}
	cfg.ServerPort = port

	cfg.DBHost = getEnvString("DB_HOST", "localhost")

	dbPort, err := getEnvInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}
	cfg.DBPort = dbPort

	debug, err := getEnvBool("DEBUG_MODE", false)
	if err != nil {
		return nil, err
	}
	cfg.DebugMode = debug

	flags := getEnvString("FEATURE_FLAGS", "")
	if flags != "" {
		flagList := strings.Split(flags, ",")
		for _, flag := range flagList {
			parts := strings.SplitN(flag, "=", 2)
			if len(parts) == 2 {
				enabled, err := strconv.ParseBool(parts[1])
				if err == nil {
					cfg.FeatureFlags[strings.TrimSpace(parts[0])] = enabled
				}
			}
		}
	}

	return cfg, nil
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	if value := os.Getenv(key); value != "" {
		return strconv.Atoi(value)
	}
	return defaultValue, nil
}

func getEnvBool(key string, defaultValue bool) (bool, error) {
	if value := os.Getenv(key); value != "" {
		return strconv.ParseBool(value)
	}
	return defaultValue, nil
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
        DatabaseURL:  getEnv("DB_URL", "postgres://localhost:5432/app"),
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
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnv(key, "")
    if val, err := strconv.ParseBool(valueStr); err == nil {
        return val
    }
    return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
    valueStr := getEnv(key, "")
    if valueStr == "" {
        return defaultValue
    }
    return strings.Split(valueStr, ",")
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

func LoadConfig() (*Config, error) {
    cfg := &Config{
        DatabaseURL:  getEnv("DB_URL", "postgres://localhost:5432/app"),
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
    if value, err := strconv.Atoi(strValue); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    strValue = strings.ToLower(strValue)
    return strValue == "true" || strValue == "1" || strValue == "yes"
}

func getEnvAsSlice(key string, defaultValue []string) []string {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    return strings.Split(strValue, ",")
}