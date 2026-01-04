package config

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
    MaxWorkers int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort: getEnvAsInt("SERVER_PORT", 8080),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnvAsInt("DB_PORT", 5432),
        DebugMode:  getEnvAsBool("DEBUG_MODE", false),
        MaxWorkers: getEnvAsInt("MAX_WORKERS", 10),
    }

    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return nil, fmt.Errorf("invalid server port: %d", cfg.ServerPort)
    }

    if cfg.MaxWorkers < 1 {
        return nil, fmt.Errorf("max workers must be positive: %d", cfg.MaxWorkers)
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
}

func (c *Config) String() string {
    var sb strings.Builder
    sb.WriteString("Configuration:\n")
    sb.WriteString(fmt.Sprintf("  ServerPort: %d\n", c.ServerPort))
    sb.WriteString(fmt.Sprintf("  DBHost: %s\n", c.DBHost))
    sb.WriteString(fmt.Sprintf("  DBPort: %d\n", c.DBPort))
    sb.WriteString(fmt.Sprintf("  DebugMode: %v\n", c.DebugMode))
    sb.WriteString(fmt.Sprintf("  MaxWorkers: %d\n", c.MaxWorkers))
    return sb.String()
}