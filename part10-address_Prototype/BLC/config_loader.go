package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST"`
		Port int    `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		Name     string `yaml:"name" env:"DB_NAME"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Format string `yaml:"format" env:"LOG_FORMAT"`
	} `yaml:"logging"`
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

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.loadFromEnv(); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) loadFromEnv() error {
	loadString := func(field *string, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			*field = val
		}
	}

	loadInt := func(field *int, envVar string) error {
		if val := os.Getenv(envVar); val != "" {
			var intVal int
			if _, err := fmt.Sscanf(val, "%d", &intVal); err != nil {
				return err
			}
			*field = intVal
		}
		return nil
	}

	loadString(&c.Server.Host, "SERVER_HOST")
	if err := loadInt(&c.Server.Port, "SERVER_PORT"); err != nil {
		return err
	}

	loadString(&c.Database.Host, "DB_HOST")
	if err := loadInt(&c.Database.Port, "DB_PORT"); err != nil {
		return err
	}
	loadString(&c.Database.Name, "DB_NAME")
	loadString(&c.Database.User, "DB_USER")
	loadString(&c.Database.Password, "DB_PASSWORD")
	loadString(&c.Database.SSLMode, "DB_SSL_MODE")

	loadString(&c.Logging.Level, "LOG_LEVEL")
	loadString(&c.Logging.Format, "LOG_FORMAT")

	return nil
}

func (c *Config) validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if c.Database.Name == "" {
		return errors.New("database name is required")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.Logging.Level] {
		return errors.New("invalid log level")
	}

	return nil
}