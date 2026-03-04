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
}