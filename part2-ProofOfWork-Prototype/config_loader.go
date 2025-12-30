package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     int    `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASS"`
	Database string `yaml:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `yaml:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Version  string         `yaml:"version"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := applyEnvOverrides(&config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func applyEnvOverrides(config *Config) error {
	envMap := map[string]string{
		"DB_HOST":               &config.Database.Host,
		"DB_PORT":               stringPtrFromInt(config.Database.Port),
		"DB_USER":               &config.Database.Username,
		"DB_PASS":               &config.Database.Password,
		"DB_NAME":               &config.Database.Database,
		"SERVER_PORT":           stringPtrFromInt(config.Server.Port),
		"SERVER_READ_TIMEOUT":   stringPtrFromInt(config.Server.ReadTimeout),
		"SERVER_WRITE_TIMEOUT":  stringPtrFromInt(config.Server.WriteTimeout),
		"DEBUG_MODE":            stringPtrFromBool(config.Server.DebugMode),
		"LOG_LEVEL":             &config.Server.LogLevel,
	}

	for envVar, fieldPtr := range envMap {
		if val := os.Getenv(envVar); val != "" {
			*fieldPtr = val
		}
	}

	return nil
}

func validateConfig(config *Config) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Server.LogLevel == "" {
		config.Server.LogLevel = "info"
	}

	return nil
}

func stringPtrFromInt(i int) *string {
	s := fmt.Sprintf("%d", i)
	return &s
}

func stringPtrFromBool(b bool) *string {
	s := fmt.Sprintf("%t", b)
	return &s
}package config

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
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	DebugMode    bool   `yaml:"debug_mode"`
	LogLevel     string `yaml:"log_level"`
}

type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(filepath string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(filepath)
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

	if config.Database.Host == "" {
		log.Print("Database host cannot be empty")
		return false
	}

	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		log.Printf("Invalid database port: %d", config.Database.Port)
		return false
	}

	return true
}