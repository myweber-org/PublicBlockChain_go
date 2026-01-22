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
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
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

	if err := overrideFromEnv(&cfg); err != nil {
		return nil, err
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func overrideFromEnv(cfg *Config) error {
	envMap := map[string]string{
		"SERVER_HOST":     &cfg.Server.Host,
		"SERVER_PORT":     stringPtrFromInt(cfg.Server.Port),
		"DB_HOST":         &cfg.Database.Host,
		"DB_PORT":         stringPtrFromInt(cfg.Database.Port),
		"DB_NAME":         &cfg.Database.Name,
		"DB_USER":         &cfg.Database.User,
		"DB_PASSWORD":     &cfg.Database.Password,
		"LOG_LEVEL":       &cfg.Logging.Level,
		"LOG_OUTPUT":      &cfg.Logging.Output,
	}

	for envVar, fieldPtr := range envMap {
		if val := os.Getenv(envVar); val != "" {
			*fieldPtr = val
		}
	}

	return nil
}

func stringPtrFromInt(i int) *string {
	s := fmt.Sprintf("%d", i)
	return &s
}

func validateConfig(cfg *Config) error {
	if cfg.Server.Host == "" {
		return errors.New("server host cannot be empty")
	}
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if cfg.Database.Name == "" {
		return errors.New("database name cannot be empty")
	}
	return nil
}package config

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
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    Debug        bool   `yaml:"debug" env:"DEBUG"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
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

    if err := overrideFromEnv(&config); err != nil {
        return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
    }

    return &config, nil
}

func overrideFromEnv(config *AppConfig) error {
    if val := os.Getenv("LOG_LEVEL"); val != "" {
        config.LogLevel = val
    }

    if val := os.Getenv("DB_HOST"); val != "" {
        config.Database.Host = val
    }
    if val := os.Getenv("DB_PORT"); val != "" {
        var port int
        if _, err := fmt.Sscanf(val, "%d", &port); err == nil {
            config.Database.Port = port
        }
    }
    if val := os.Getenv("DB_USER"); val != "" {
        config.Database.Username = val
    }
    if val := os.Getenv("DB_PASS"); val != "" {
        config.Database.Password = val
    }
    if val := os.Getenv("DB_NAME"); val != "" {
        config.Database.Name = val
    }

    if val := os.Getenv("SERVER_PORT"); val != "" {
        var port int
        if _, err := fmt.Sscanf(val, "%d", &port); err == nil {
            config.Server.Port = port
        }
    }
    if val := os.Getenv("READ_TIMEOUT"); val != "" {
        var timeout int
        if _, err := fmt.Sscanf(val, "%d", &timeout); err == nil {
            config.Server.ReadTimeout = timeout
        }
    }
    if val := os.Getenv("WRITE_TIMEOUT"); val != "" {
        var timeout int
        if _, err := fmt.Sscanf(val, "%d", &timeout); err == nil {
            config.Server.WriteTimeout = timeout
        }
    }
    if val := os.Getenv("DEBUG"); val != "" {
        config.Server.Debug = (val == "true" || val == "1")
    }

    return nil
}