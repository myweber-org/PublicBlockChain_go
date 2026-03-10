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
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    if configPath == "" {
        configPath = "config.yaml"
    }

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to get absolute path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    overrideString(&config.Database.Host, "DB_HOST")
    overrideString(&config.Database.Username, "DB_USER")
    overrideString(&config.Database.Password, "DB_PASS")
    overrideString(&config.Database.Name, "DB_NAME")
    overrideInt(&config.Database.Port, "DB_PORT")
    
    overrideInt(&config.Server.Port, "SERVER_PORT")
    overrideInt(&config.Server.ReadTimeout, "READ_TIMEOUT")
    overrideInt(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
    overrideBool(&config.Server.DebugMode, "DEBUG_MODE")
    
    overrideString(&config.LogLevel, "LOG_LEVEL")
}

func overrideString(field *string, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        *field = val
    }
}

func overrideInt(field *int, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        var intVal int
        if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
            *field = intVal
        }
    }
}

func overrideBool(field *bool, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        *field = val == "true" || val == "1" || val == "yes"
    }
}package config

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
    LogLevel string
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{
        ServerPort:     8080,
        DatabaseURL:    "localhost:5432",
        CacheEnabled:   true,
        MaxConnections: 100,
        LogLevel:       "info",
    }

    if portStr := os.Getenv("APP_PORT"); portStr != "" {
        port, err := strconv.Atoi(portStr)
        if err != nil {
            return nil, fmt.Errorf("invalid port value: %v", err)
        }
        cfg.ServerPort = port
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if cacheStr := os.Getenv("CACHE_ENABLED"); cacheStr != "" {
        cacheEnabled, err := strconv.ParseBool(cacheStr)
        if err != nil {
            return nil, fmt.Errorf("invalid cache enabled value: %v", err)
        }
        cfg.CacheEnabled = cacheEnabled
    }

    if maxConnStr := os.Getenv("MAX_CONNECTIONS"); maxConnStr != "" {
        maxConn, err := strconv.Atoi(maxConnStr)
        if err != nil {
            return nil, fmt.Errorf("invalid max connections value: %v", err)
        }
        cfg.MaxConnections = maxConn
    }

    if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
        validLevels := map[string]bool{
            "debug": true,
            "info":  true,
            "warn":  true,
            "error": true,
        }
        if !validLevels[strings.ToLower(logLevel)] {
            return nil, fmt.Errorf("invalid log level: %s", logLevel)
        }
        cfg.LogLevel = strings.ToLower(logLevel)
    }

    return cfg, nil
}

func ValidateConfig(cfg *AppConfig) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }

    if cfg.DatabaseURL == "" {
        return fmt.Errorf("database URL cannot be empty")
    }

    if cfg.MaxConnections < 1 {
        return fmt.Errorf("max connections must be greater than 0")
    }

    return nil
}