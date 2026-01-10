package config

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

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Version  string         `yaml:"version"`
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
    if envVal := os.Getenv("DB_HOST"); envVal != "" {
        config.Database.Host = envVal
    }
    if envVal := os.Getenv("DB_PORT"); envVal != "" {
        var port int
        if _, err := fmt.Sscanf(envVal, "%d", &port); err == nil {
            config.Database.Port = port
        }
    }
    if envVal := os.Getenv("DB_USER"); envVal != "" {
        config.Database.Username = envVal
    }
    if envVal := os.Getenv("DB_PASS"); envVal != "" {
        config.Database.Password = envVal
    }
    if envVal := os.Getenv("DB_NAME"); envVal != "" {
        config.Database.Name = envVal
    }
    if envVal := os.Getenv("SERVER_PORT"); envVal != "" {
        var port int
        if _, err := fmt.Sscanf(envVal, "%d", &port); err == nil {
            config.Server.Port = port
        }
    }
    if envVal := os.Getenv("READ_TIMEOUT"); envVal != "" {
        var timeout int
        if _, err := fmt.Sscanf(envVal, "%d", &timeout); err == nil {
            config.Server.ReadTimeout = timeout
        }
    }
    if envVal := os.Getenv("WRITE_TIMEOUT"); envVal != "" {
        var timeout int
        if _, err := fmt.Sscanf(envVal, "%d", &timeout); err == nil {
            config.Server.WriteTimeout = timeout
        }
    }
    if envVal := os.Getenv("DEBUG_MODE"); envVal != "" {
        config.Server.DebugMode = envVal == "true" || envVal == "1"
    }
}

func (c *AppConfig) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if c.Database.Port <= 0 || c.Database.Port > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    return nil
}