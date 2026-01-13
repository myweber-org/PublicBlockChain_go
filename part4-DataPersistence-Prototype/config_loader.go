package config

import (
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
	cfg := &Config{}

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

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	overrideFromEnv(cfg)

	return cfg, nil
}

func overrideFromEnv(cfg *Config) {
	if val := os.Getenv("SERVER_HOST"); val != "" {
		cfg.Server.Host = val
	}
	if val := os.Getenv("SERVER_PORT"); val != "" {
		port := 0
		fmt.Sscanf(val, "%d", &port)
		if port > 0 {
			cfg.Server.Port = port
		}
	}
	if val := os.Getenv("DB_HOST"); val != "" {
		cfg.Database.Host = val
	}
	if val := os.Getenv("DB_PORT"); val != "" {
		port := 0
		fmt.Sscanf(val, "%d", &port)
		if port > 0 {
			cfg.Database.Port = port
		}
	}
	if val := os.Getenv("DB_NAME"); val != "" {
		cfg.Database.Name = val
	}
	if val := os.Getenv("DB_USER"); val != "" {
		cfg.Database.User = val
	}
	if val := os.Getenv("DB_PASSWORD"); val != "" {
		cfg.Database.Password = val
	}
	if val := os.Getenv("DB_SSL_MODE"); val != "" {
		cfg.Database.SSLMode = val
	}
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		cfg.Logging.Level = val
	}
	if val := os.Getenv("LOG_FORMAT"); val != "" {
		cfg.Logging.Format = val
	}
}