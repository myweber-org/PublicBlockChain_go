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
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Format string `yaml:"format" env:"LOG_FORMAT"`
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
	cfg.Server.Host = getEnvOrDefault("SERVER_HOST", cfg.Server.Host)
	cfg.Server.Port = getEnvIntOrDefault("SERVER_PORT", cfg.Server.Port)
	cfg.Database.Host = getEnvOrDefault("DB_HOST", cfg.Database.Host)
	cfg.Database.Port = getEnvIntOrDefault("DB_PORT", cfg.Database.Port)
	cfg.Database.Name = getEnvOrDefault("DB_NAME", cfg.Database.Name)
	cfg.Database.User = getEnvOrDefault("DB_USER", cfg.Database.User)
	cfg.Database.Password = getEnvOrDefault("DB_PASSWORD", cfg.Database.Password)
	cfg.Logging.Level = getEnvOrDefault("LOG_LEVEL", cfg.Logging.Level)
	cfg.Logging.Format = getEnvOrDefault("LOG_FORMAT", cfg.Logging.Format)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}package config

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type ServerConfig struct {
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	Environment  string `yaml:"environment"`
}

type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Features []string       `yaml:"features"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
	if filePath == "" {
		return nil, errors.New("config file path cannot be empty")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 {
		return errors.New("database port must be positive")
	}
	if config.Server.Port <= 0 {
		return errors.New("server port must be positive")
	}
	if config.Server.Environment == "" {
		config.Server.Environment = "development"
	}
	return nil
}package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	ServerPort int    `json:"server_port" env:"SERVER_PORT"`
	DBHost     string `json:"db_host" env:"DB_HOST"`
	DBPort     int    `json:"db_port" env:"DB_PORT"`
	DebugMode  bool   `json:"debug_mode" env:"DEBUG_MODE"`
	APIKey     string `json:"api_key" env:"API_KEY"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{
		ServerPort: 8080,
		DBHost:     "localhost",
		DBPort:     5432,
		DebugMode:  false,
	}

	if configPath != "" {
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			return nil, fmt.Errorf("invalid config path: %w", err)
		}

		data, err := os.ReadFile(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config JSON: %w", err)
		}
	}

	if err := loadEnvVars(cfg); err != nil {
		return nil, err
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadEnvVars(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}

		fieldValue := v.Field(i)
		switch field.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(envValue)
		case reflect.Int:
			var intVal int
			if _, err := fmt.Sscanf(envValue, "%d", &intVal); err != nil {
				return fmt.Errorf("invalid integer value for %s: %s", envTag, envValue)
			}
			fieldValue.SetInt(int64(intVal))
		case reflect.Bool:
			boolVal := strings.ToLower(envValue) == "true" || envValue == "1"
			fieldValue.SetBool(boolVal)
		}
	}

	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.ServerPort)
	}

	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return fmt.Errorf("invalid database port: %d", cfg.DBPort)
	}

	if cfg.DBHost == "" {
		return fmt.Errorf("database host cannot be empty")
	}

	return nil
}package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    MaxConnections int
    AllowedOrigins []string
}

func Load() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnv("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/app")
    
    cacheStr := getEnv("CACHE_ENABLED", "true")
    cfg.CacheEnabled = strings.ToLower(cacheStr) == "true"
    
    maxConnStr := getEnv("MAX_CONNECTIONS", "100")
    maxConn, err := strconv.Atoi(maxConnStr)
    if err != nil {
        return nil, err
    }
    cfg.MaxConnections = maxConn
    
    originsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080")
    cfg.AllowedOrigins = strings.Split(originsStr, ",")
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}