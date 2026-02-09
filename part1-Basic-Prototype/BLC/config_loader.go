package config

import (
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v2"
)

type AppConfig struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Name     string `yaml:"name"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
    } `yaml:"database"`
    Logging struct {
        Level  string `yaml:"level"`
        Output string `yaml:"output"`
    } `yaml:"logging"`
}

func LoadConfig(path string) (*AppConfig, error) {
    data, err := ioutil.ReadFile(path)
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
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        log.Printf("Invalid server port: %d", config.Server.Port)
        return false
    }

    if config.Database.Host == "" || config.Database.Name == "" {
        log.Print("Database host and name must be specified")
        return false
    }

    if config.Logging.Level != "debug" && 
       config.Logging.Level != "info" && 
       config.Logging.Level != "warn" && 
       config.Logging.Level != "error" {
        log.Printf("Invalid logging level: %s", config.Logging.Level)
        return false
    }

    return true
}