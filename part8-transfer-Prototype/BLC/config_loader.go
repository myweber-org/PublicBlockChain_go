package config

import (
    "fmt"
    "io"
    "os"

    "gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port         int    `yaml:"port"`
    ReadTimeout  int    `yaml:"read_timeout"`
    WriteTimeout int    `yaml:"write_timeout"`
}

type AppConfig struct {
    Environment string        `yaml:"environment"`
    Database    DatabaseConfig `yaml:"database"`
    Server      ServerConfig   `yaml:"server"`
}

func LoadConfig(path string) (*AppConfig, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(config *AppConfig) error {
    if config.Environment == "" {
        return fmt.Errorf("environment must be specified")
    }

    if config.Database.Host == "" {
        return fmt.Errorf("database host must be specified")
    }

    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }

    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }

    return nil
}