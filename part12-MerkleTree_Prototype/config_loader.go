package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port string `yaml:"port" env:"SERVER_PORT"`
		Host string `yaml:"host" env:"SERVER_HOST"`
	} `yaml:"server"`
	Database struct {
		URL      string `yaml:"url" env:"DB_URL"`
		MaxConns int    `yaml:"max_connections" env:"DB_MAX_CONNS"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(filePath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	overrideWithEnv(config)

	return config, nil
}

func overrideWithEnv(c *Config) {
	envOverride(&c.Server.Port, "SERVER_PORT")
	envOverride(&c.Server.Host, "SERVER_HOST")
	envOverride(&c.Database.URL, "DB_URL")
	envOverrideInt(&c.Database.MaxConns, "DB_MAX_CONNS")
	envOverride(&c.LogLevel, "LOG_LEVEL")
}

func envOverride(field *string, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val
	}
}

func envOverrideInt(field *int, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		if intVal, err := parseInt(val); err == nil {
			*field = intVal
		}
	}
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}