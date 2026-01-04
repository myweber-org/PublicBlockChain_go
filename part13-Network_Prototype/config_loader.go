package config

import (
	"errors"
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

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	var err error

	cfg.ServerPort, err = getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}

	cfg.DBHost = getEnvString("DB_HOST", "localhost")
	
	cfg.DBPort, err = getEnvInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	cfg.DebugMode, err = getEnvBool("DEBUG_MODE", false)
	if err != nil {
		return nil, err
	}

	apiKeysStr := getEnvString("API_KEYS", "")
	if apiKeysStr != "" {
		cfg.APIKeys = strings.Split(apiKeysStr, ",")
		for i, key := range cfg.APIKeys {
			cfg.APIKeys[i] = strings.TrimSpace(key)
		}
	}

	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return nil, errors.New("invalid server port range")
	}

	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return nil, errors.New("invalid database port range")
	}

	return cfg, nil
}

func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("invalid integer value for " + key)
		}
		return intValue, nil
	}
	return defaultValue, nil
}

func getEnvBool(key string, defaultValue bool) (bool, error) {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return false, errors.New("invalid boolean value for " + key)
		}
		return boolValue, nil
	}
	return defaultValue, nil
}