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
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
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

	return &cfg, nil
}

func overrideFromEnv(cfg *Config) error {
	envMap := map[string]*string{
		"SERVER_HOST":     &cfg.Server.Host,
		"DB_HOST":         &cfg.Database.Host,
		"DB_NAME":         &cfg.Database.Name,
		"DB_USER":         &cfg.Database.User,
		"DB_PASSWORD":     &cfg.Database.Password,
		"LOG_LEVEL":       &cfg.LogLevel,
	}

	for envVar, field := range envMap {
		if val, exists := os.LookupEnv(envVar); exists && val != "" {
			*field = val
		}
	}

	portEnvVars := map[string]*int{
		"SERVER_PORT": &cfg.Server.Port,
		"DB_PORT":     &cfg.Database.Port,
	}

	for envVar, field := range portEnvVars {
		if val, exists := os.LookupEnv(envVar); exists && val != "" {
			port, err := parseInt(val)
			if err != nil {
				return err
			}
			*field = port
		}
	}

	return nil
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0, errors.New("invalid integer value")
	}
	return result, nil
}package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DBHost     string
    DBPort     int
    DebugMode  bool
    APIKeys    []string
}

func Load() (*Config, error) {
    cfg := &Config{}
    var err error

    cfg.ServerPort, err = getIntEnv("SERVER_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }

    cfg.DBHost = getStringEnv("DB_HOST", "localhost")
    
    cfg.DBPort, err = getIntEnv("DB_PORT", 5432)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }

    cfg.DebugMode, err = getBoolEnv("DEBUG_MODE", false)
    if err != nil {
        return nil, fmt.Errorf("invalid DEBUG_MODE: %w", err)
    }

    cfg.APIKeys = getSliceEnv("API_KEYS", []string{"default_key"})

    return cfg, nil
}

func getStringEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
    if value := os.Getenv(key); value != "" {
        return strconv.Atoi(value)
    }
    return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) (bool, error) {
    if value := os.Getenv(key); value != "" {
        return strconv.ParseBool(value)
    }
    return defaultValue, nil
}

func getSliceEnv(key string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        return strings.Split(value, ",")
    }
    return defaultValue
}