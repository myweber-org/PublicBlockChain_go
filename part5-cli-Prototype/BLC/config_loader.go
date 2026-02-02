package config

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
    Debug        bool   `yaml:"debug" env:"SERVER_DEBUG"`
    LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Version  string         `yaml:"version"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    if configPath == "" {
        configPath = "config.yaml"
    }

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to get absolute path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    overrideString(&config.Database.Host, "DB_HOST")
    overrideInt(&config.Database.Port, "DB_PORT")
    overrideString(&config.Database.Username, "DB_USER")
    overrideString(&config.Database.Password, "DB_PASS")
    overrideString(&config.Database.Name, "DB_NAME")
    
    overrideInt(&config.Server.Port, "SERVER_PORT")
    overrideBool(&config.Server.Debug, "SERVER_DEBUG")
    overrideString(&config.Server.LogLevel, "LOG_LEVEL")
    overrideInt(&config.Server.ReadTimeout, "READ_TIMEOUT")
    overrideInt(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
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
}

func overrideBool(field *bool, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        *field = val == "true" || val == "1" || val == "yes"
    }
}package config

import (
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
    LogLevel string `yaml:"log_level"`
}

func LoadConfig(path string) (*Config, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("config file does not exist: %s", path)
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %v", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %v", err)
    }

    if err := validateConfig(&cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %v", err)
    }

    return &cfg, nil
}

func validateConfig(cfg *Config) error {
    if cfg.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if cfg.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    if cfg.LogLevel == "" {
        cfg.LogLevel = "info"
    }

    return nil
}package config

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level"`
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path cannot be empty")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(c *Config) error {
	if c.Server.Host == "" {
		return errors.New("server host is required")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if c.Database.Host == "" {
		return errors.New("database host is required")
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
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
    ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug_mode" env:"SERVER_DEBUG"`
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

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    config.Database.Host = getEnvOrDefault("DB_HOST", config.Database.Host)
    config.Database.Port = getEnvIntOrDefault("DB_PORT", config.Database.Port)
    config.Database.Username = getEnvOrDefault("DB_USER", config.Database.Username)
    config.Database.Password = getEnvOrDefault("DB_PASS", config.Database.Password)
    config.Database.Name = getEnvOrDefault("DB_NAME", config.Database.Name)

    config.Server.Port = getEnvIntOrDefault("SERVER_PORT", config.Server.Port)
    config.Server.ReadTimeout = getEnvIntOrDefault("SERVER_READ_TIMEOUT", config.Server.ReadTimeout)
    config.Server.WriteTimeout = getEnvIntOrDefault("SERVER_WRITE_TIMEOUT", config.Server.WriteTimeout)
    config.Server.DebugMode = getEnvBoolOrDefault("SERVER_DEBUG", config.Server.DebugMode)

    config.LogLevel = getEnvOrDefault("LOG_LEVEL", config.LogLevel)
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        var intValue int
        if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        return value == "true" || value == "1" || value == "yes"
    }
    return defaultValue
}