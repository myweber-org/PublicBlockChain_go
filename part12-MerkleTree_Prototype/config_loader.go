package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	AllowedHosts []string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	
	portStr := getEnv("SERVER_PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	cfg.ServerPort = port
	
	debugStr := getEnv("DEBUG_MODE", "false")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"
	
	cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/appdb")
	
	hostsStr := getEnv("ALLOWED_HOSTS", "localhost,127.0.0.1")
	cfg.AllowedHosts = strings.Split(hostsStr, ",")
	
	return cfg, nil
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
	DebugMode  bool
	MaxWorkers int
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerPort: 8080,
		DebugMode:  false,
		MaxWorkers: 10,
	}

	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.ServerPort = port
		}
	}

	if debugStr := os.Getenv("DEBUG_MODE"); debugStr != "" {
		if debug, err := strconv.ParseBool(debugStr); err == nil {
			cfg.DebugMode = debug
		}
	}

	if workersStr := os.Getenv("MAX_WORKERS"); workersStr != "" {
		if workers, err := strconv.Atoi(workersStr); err == nil {
			cfg.MaxWorkers = workers
		}
	}

	return cfg, nil
}