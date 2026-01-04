package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
    Name     string `yaml:"name" env:"DB_NAME"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
    Debug        bool   `yaml:"debug" env:"SERVER_DEBUG"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    var config AppConfig

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    config.Database.Host = getEnvOrDefault("DB_HOST", config.Database.Host)
    config.Database.Port = getEnvIntOrDefault("DB_PORT", config.Database.Port)
    config.Database.Username = getEnvOrDefault("DB_USER", config.Database.Username)
    config.Database.Password = getEnvOrDefault("DB_PASS", config.Database.Password)
    config.Database.Name = getEnvOrDefault("DB_NAME", config.Database.Name)

    config.Server.Port = getEnvIntOrDefault("SERVER_PORT", config.Server.Port)
    config.Server.ReadTimeout = getEnvIntOrDefault("SERVER_READ_TIMEOUT", config.Server.ReadTimeout)
    config.Server.WriteTimeout = getEnvIntOrDefault("SERVER_WRITE_TIMEOUT", config.Server.WriteTimeout)
    config.Server.Debug = getEnvBoolOrDefault("SERVER_DEBUG", config.Server.Debug)

    config.LogLevel = getEnvOrDefault("LOG_LEVEL", config.LogLevel)
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        var result int
        if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
            return result
        }
    }
    return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        return value == "true" || value == "1" || value == "yes"
    }
    return defaultValue
}