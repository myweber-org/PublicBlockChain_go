package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DebugMode  bool   `json:"debug_mode"`
}

var (
	instance *Config
	once     sync.Once
)

func Load() *Config {
	once.Do(func() {
		instance = &Config{
			ServerPort: getEnv("SERVER_PORT", "8080"),
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DebugMode:  getEnv("DEBUG_MODE", "false") == "true",
		}

		if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
			loadFromFile(configFile, instance)
		}
	})
	return instance
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadFromFile(filename string, config *Config) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	var fileConfig Config
	if err := json.Unmarshal(data, &fileConfig); err != nil {
		return
	}

	if fileConfig.ServerPort != "" {
		config.ServerPort = fileConfig.ServerPort
	}
	if fileConfig.DBHost != "" {
		config.DBHost = fileConfig.DBHost
	}
	if fileConfig.DBPort != "" {
		config.DBPort = fileConfig.DBPort
	}
	config.DebugMode = config.DebugMode || fileConfig.DebugMode
}