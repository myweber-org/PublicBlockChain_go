package config

import (
    "io/ioutil"
    "log"

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
    Port         int            `yaml:"port"`
    Debug        bool           `yaml:"debug"`
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var config AppConfig
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }

    return &config, nil
}

func ValidateConfig(config *AppConfig) bool {
    if config.Database.Host == "" {
        log.Println("Database host is required")
        return false
    }
    if config.Database.Port <= 0 {
        log.Println("Database port must be positive")
        return false
    }
    if config.Server.Port <= 0 {
        log.Println("Server port must be positive")
        return false
    }
    return true
}