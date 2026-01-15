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
}package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v2"
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
    ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug_mode" env:"SERVER_DEBUG"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Version  string         `yaml:"version"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    var config AppConfig

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    loadFromEnv(&config.Database)
    loadFromEnv(&config.Server)

    return &config, nil
}

func loadFromEnv(config interface{}) {
    // Environment variable loading implementation
    // This would use reflection to check struct tags
    // and override values from environment variables
}