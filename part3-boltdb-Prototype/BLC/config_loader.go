package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    DatabaseURL string
    MaxConnections int
    DebugMode bool
    AllowedOrigins []string
}

func Load() (*Config, error) {
    cfg := &Config{
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
        DebugMode: getEnvAsBool("DEBUG_MODE", false),
        AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"localhost:3000"}, ","),
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

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
    valueStr := getEnv(key, "")
    if valueStr == "" {
        return defaultValue
    }
    return strings.Split(valueStr, sep)
}package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    DatabaseURL string
    MaxConnections int
    DebugMode bool
    AllowedOrigins []string
}

func Load() (*Config, error) {
    cfg := &Config{
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
        DebugMode: getEnvAsBool("DEBUG_MODE", false),
        AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"localhost:3000"}),
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

func getEnvAsSlice(key string, defaultValue []string) []string {
    valueStr := getEnv(key, "")
    if valueStr == "" {
        return defaultValue
    }
    return strings.Split(valueStr, ",")
}package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    EnableDebug bool
    AllowedOrigins []string
}

func Load() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnv("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/app")
    
    debugStr := getEnv("ENABLE_DEBUG", "false")
    cfg.EnableDebug = strings.ToLower(debugStr) == "true"
    
    originsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
    cfg.AllowedOrigins = strings.Split(originsStr, ",")
    
    if err := validateConfig(cfg); err != nil {
        return nil, err
    }
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return strconv.ErrRange
    }
    
    if cfg.DatabaseURL == "" {
        return strconv.ErrSyntax
    }
    
    return nil
}package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v2"
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

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
    var cfg Config

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    overrideFromEnv(&cfg)

    return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
    if val := os.Getenv("DB_HOST"); val != "" {
        cfg.Database.Host = val
    }
    if val := os.Getenv("DB_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &cfg.Database.Port)
    }
    if val := os.Getenv("DB_USER"); val != "" {
        cfg.Database.Username = val
    }
    if val := os.Getenv("DB_PASS"); val != "" {
        cfg.Database.Password = val
    }
    if val := os.Getenv("DB_NAME"); val != "" {
        cfg.Database.Name = val
    }
    if val := os.Getenv("SERVER_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &cfg.Server.Port)
    }
    if val := os.Getenv("READ_TIMEOUT"); val != "" {
        fmt.Sscanf(val, "%d", &cfg.Server.ReadTimeout)
    }
    if val := os.Getenv("WRITE_TIMEOUT"); val != "" {
        fmt.Sscanf(val, "%d", &cfg.Server.WriteTimeout)
    }
    if val := os.Getenv("DEBUG_MODE"); val != "" {
        cfg.Server.DebugMode = (val == "true" || val == "1")
    }
    if val := os.Getenv("LOG_LEVEL"); val != "" {
        cfg.LogLevel = val
    }
}

func (c *Config) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if c.Database.Port <= 0 || c.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", c.Database.Port)
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }
    return nil
}