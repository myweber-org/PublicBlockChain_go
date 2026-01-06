
package config

import (
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	AllowedHosts []string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	
	portStr := getEnvWithDefault("SERVER_PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	cfg.ServerPort = port
	
	debugStr := getEnvWithDefault("DEBUG_MODE", "false")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"
	
	cfg.DatabaseURL = getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/appdb")
	
	hostsStr := getEnvWithDefault("ALLOWED_HOSTS", "localhost,127.0.0.1")
	cfg.AllowedHosts = strings.Split(hostsStr, ",")
	
	return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}