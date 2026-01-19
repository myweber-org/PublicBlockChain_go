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
    Debug        bool   `yaml:"debug" env:"SERVER_DEBUG"`
    LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Features []string       `yaml:"features"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    data, err := os.ReadFile(configPath)
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
    overrideInt(&config.Database.Port, "DB_PORT")
    overrideString(&config.Database.Username, "DB_USER")
    overrideString(&config.Database.Password, "DB_PASS")
    overrideString(&config.Database.Name, "DB_NAME")
    
    overrideInt(&config.Server.Port, "SERVER_PORT")
    overrideBool(&config.Server.Debug, "SERVER_DEBUG")
    overrideString(&config.Server.LogLevel, "LOG_LEVEL")
    overrideInt(&config.Server.ReadTimeout, "READ_TIMEOUT")
    overrideInt(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
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
}

func DefaultConfigPath() string {
    paths := []string{
        "config.yaml",
        "config.yml",
        filepath.Join("config", "app.yaml"),
        filepath.Join("etc", "app", "config.yaml"),
    }
    
    for _, path := range paths {
        if _, err := os.Stat(path); err == nil {
            return path
        }
    }
    return ""
}