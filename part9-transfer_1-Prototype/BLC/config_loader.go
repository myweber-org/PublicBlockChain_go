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
	envMap := map[string]*string{
		"SERVER_HOST":     &cfg.Server.Host,
		"DB_HOST":         &cfg.Database.Host,
		"DB_NAME":         &cfg.Database.Name,
		"DB_USER":         &cfg.Database.User,
		"DB_PASSWORD":     &cfg.Database.Password,
		"LOG_LEVEL":       &cfg.Logging.Level,
		"LOG_OUTPUT":      &cfg.Logging.Output,
	}

	for envVar, field := range envMap {
		if val, exists := os.LookupEnv(envVar); exists && val != "" {
			*field = val
		}
	}

	if portStr, exists := os.LookupEnv("SERVER_PORT"); exists && portStr != "" {
		if port, err := parseInt(portStr); err == nil {
			cfg.Server.Port = port
		}
	}

	if portStr, exists := os.LookupEnv("DB_PORT"); exists && portStr != "" {
		if port, err := parseInt(portStr); err == nil {
			cfg.Database.Port = port
		}
	}

	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.Server.Host == "" {
		return errors.New("server host is required")
	}
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return errors.New("invalid server port")
	}
	if cfg.Database.Host == "" {
		return errors.New("database host is required")
	}
	if cfg.Database.Name == "" {
		return errors.New("database name is required")
	}
	return nil
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}package config

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
	} `yaml:"database"`
	Logging struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
		File  string `yaml:"file" env:"LOG_FILE"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
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

	overrideFromEnv(&cfg)
	return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
	overrideString(&cfg.Server.Host, "SERVER_HOST")
	overrideInt(&cfg.Server.Port, "SERVER_PORT")
	overrideString(&cfg.Database.Host, "DB_HOST")
	overrideInt(&cfg.Database.Port, "DB_PORT")
	overrideString(&cfg.Database.Name, "DB_NAME")
	overrideString(&cfg.Database.User, "DB_USER")
	overrideString(&cfg.Database.Password, "DB_PASSWORD")
	overrideString(&cfg.Logging.Level, "LOG_LEVEL")
	overrideString(&cfg.Logging.File, "LOG_FILE")
}

func overrideString(field *string, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val
	}
}

func overrideInt(field *int, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		var intVal int
		if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
			*field = intVal
		}
	}
}package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DebugMode  bool
    DatabaseURL string
    AllowedHosts []string
}

func Load() (*Config, error) {
    cfg := &Config{
        ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
        DebugMode:   getEnvAsBool("DEBUG_MODE", false),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost"}, ","),
    }
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnv(key, "")
    if val, err := strconv.ParseBool(valueStr); err == nil {
        return val
    }
    return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
    valueStr := getEnv(key, "")
    if valueStr == "" {
        return defaultValue
    }
    return strings.Split(valueStr, sep)
}