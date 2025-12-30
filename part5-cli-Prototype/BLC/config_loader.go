package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    MaxConnections int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvOrDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnvOrDefault("DATABASE_URL", "postgres://localhost:5432/app")
    
    cacheStr := getEnvOrDefault("CACHE_ENABLED", "true")
    cfg.CacheEnabled = strings.ToLower(cacheStr) == "true"
    
    maxConnStr := getEnvOrDefault("MAX_CONNECTIONS", "100")
    maxConn, err := strconv.Atoi(maxConnStr)
    if err != nil {
        return nil, err
    }
    cfg.MaxConnections = maxConn
    
    return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}