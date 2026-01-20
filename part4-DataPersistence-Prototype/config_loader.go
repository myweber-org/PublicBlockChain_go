package config

import (
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
    Logging struct {
        Level string `yaml:"level"`
        File  string `yaml:"file"`
    } `yaml:"logging"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(c *Config) error {
    if c.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if c.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    if c.Logging.Level == "" {
        c.Logging.Level = "info"
    }
    return nil
}

func SaveConfig(c *Config, path string) error {
    data, err := yaml.Marshal(c)
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }

    if err := ioutil.WriteFile(path, data, 0644); err != nil {
        return fmt.Errorf("failed to write config file: %w", err)
    }

    return nil
}

func GenerateDefaultConfig() *Config {
    var config Config
    config.Server.Host = "localhost"
    config.Server.Port = 8080
    config.Database.Host = "localhost"
    config.Database.Port = 5432
    config.Database.Username = "postgres"
    config.Database.Name = "appdb"
    config.Logging.Level = "info"
    config.Logging.File = "app.log"
    return &config
}

func ConfigExists(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}