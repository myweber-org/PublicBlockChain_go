package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    MaxConnections int
    FeatureFlags map[string]bool
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{
        FeatureFlags: make(map[string]bool),
    }

    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT value: %v", err)
    }
    cfg.ServerPort = port

    cfg.DatabaseURL = getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/appdb")

    cacheEnabledStr := getEnvWithDefault("CACHE_ENABLED", "true")
    cfg.CacheEnabled = strings.ToLower(cacheEnabledStr) == "true"

    maxConnStr := getEnvWithDefault("MAX_CONNECTIONS", "100")
    maxConn, err := strconv.Atoi(maxConnStr)
    if err != nil {
        return nil, fmt.Errorf("invalid MAX_CONNECTIONS value: %v", err)
    }
    cfg.MaxConnections = maxConn

    featureFlags := getEnvWithDefault("FEATURE_FLAGS", "new_ui=false,experimental_api=true")
    parseFeatureFlags(featureFlags, cfg.FeatureFlags)

    if err := validateConfig(cfg); err != nil {
        return nil, err
    }

    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func parseFeatureFlags(flagsStr string, flagsMap map[string]bool) {
    pairs := strings.Split(flagsStr, ",")
    for _, pair := range pairs {
        parts := strings.Split(pair, "=")
        if len(parts) == 2 {
            key := strings.TrimSpace(parts[0])
            value := strings.TrimSpace(parts[1])
            flagsMap[key] = strings.ToLower(value) == "true"
        }
    }
}

func validateConfig(cfg *AppConfig) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if cfg.MaxConnections < 1 {
        return fmt.Errorf("max connections must be positive")
    }
    if cfg.DatabaseURL == "" {
        return fmt.Errorf("database URL cannot be empty")
    }
    return nil
}