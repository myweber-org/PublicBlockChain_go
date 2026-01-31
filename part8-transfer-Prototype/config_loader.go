package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	ServerPort string `json:"server_port" env:"SERVER_PORT"`
	DBHost     string `json:"db_host" env:"DB_HOST"`
	DBPort     int    `json:"db_port" env:"DB_PORT"`
	DebugMode  bool   `json:"debug_mode" env:"DEBUG_MODE"`
}

func LoadConfig(configPath string) (*Config, error) {
	var cfg Config

	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	}

	cfg.loadFromEnv()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) loadFromEnv() {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		c.ServerPort = port
	}
	if host := os.Getenv("DB_HOST"); host != "" {
		c.DBHost = host
	}
}

func (c *Config) validate() error {
	if c.ServerPort == "" {
		return errors.New("server port is required")
	}
	if c.DBHost == "" {
		return errors.New("database host is required")
	}
	if c.DBPort <= 0 || c.DBPort > 65535 {
		return errors.New("invalid database port")
	}
	return nil
}

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}