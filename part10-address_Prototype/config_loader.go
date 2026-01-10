
package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    LogLevel string
    MaxConnections int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvOrDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnvOrDefault("DATABASE_URL", "postgres://localhost:5432/app")
    
    cacheEnabledStr := getEnvOrDefault("CACHE_ENABLED", "true")
    cfg.CacheEnabled = strings.ToLower(cacheEnabledStr) == "true"
    
    cfg.LogLevel = strings.ToUpper(getEnvOrDefault("LOG_LEVEL", "INFO"))
    validLogLevels := map[string]bool{"DEBUG": true, "INFO": true, "WARN": true, "ERROR": true}
    if !validLogLevels[cfg.LogLevel] {
        return nil, fmt.Errorf("invalid LOG_LEVEL: %s", cfg.LogLevel)
    }
    
    maxConnStr := getEnvOrDefault("MAX_CONNECTIONS", "100")
    maxConn, err := strconv.Atoi(maxConnStr)
    if err != nil {
        return nil, fmt.Errorf("invalid MAX_CONNECTIONS: %v", err)
    }
    if maxConn <= 0 {
        return nil, fmt.Errorf("MAX_CONNECTIONS must be positive")
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