package config

import (
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port int    `yaml:"port"`
    Mode string `yaml:"mode"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return &config, nil
}

func ValidateConfig(config *AppConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Port <= 0 {
        return fmt.Errorf("database port must be positive")
    }
    if config.Server.Port <= 0 {
        return fmt.Errorf("server port must be positive")
    }
    return nil
}

func GetDefaultConfig() *AppConfig {
    return &AppConfig{
        Database: DatabaseConfig{
            Host:     "localhost",
            Port:     5432,
            Username: "postgres",
            Password: "",
            Name:     "appdb",
        },
        Server: ServerConfig{
            Port: 8080,
            Mode: "development",
        },
        LogLevel: "info",
    }
}

func SaveConfig(config *AppConfig, filePath string) error {
    data, err := yaml.Marshal(config)
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }

    if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
        return fmt.Errorf("failed to write config file: %w", err)
    }

    return nil
}

func LoadOrCreateConfig(filePath string) (*AppConfig, error) {
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        defaultConfig := GetDefaultConfig()
        if err := SaveConfig(defaultConfig, filePath); err != nil {
            return nil, fmt.Errorf("failed to create default config: %w", err)
        }
        return defaultConfig, nil
    }

    return LoadConfig(filePath)
}