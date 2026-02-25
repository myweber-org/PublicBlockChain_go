
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	APIKeys    []string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}
	if port < 1 || port > 65535 {
		return nil, errors.New("SERVER_PORT out of valid range")
	}
	cfg.ServerPort = port

	debugStr := os.Getenv("DEBUG_MODE")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	cfg.DatabaseURL = dbURL

	keysStr := os.Getenv("API_KEYS")
	if keysStr != "" {
		cfg.APIKeys = strings.Split(keysStr, ",")
		for i, key := range cfg.APIKeys {
			cfg.APIKeys[i] = strings.TrimSpace(key)
		}
	}

	return cfg, nil
}