package config

import (
    "os"
    "strconv"
    "strings"
)

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
    SSLMode  string
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
    DebugMode    bool
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    LogLevel string
}

func LoadConfig() (*Config, error) {
    dbConfig := DatabaseConfig{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnvAsInt("DB_PORT", 5432),
        Username: getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASSWORD", ""),
        Database: getEnv("DB_NAME", "appdb"),
        SSLMode:  getEnv("DB_SSL_MODE", "disable"),
    }

    serverConfig := ServerConfig{
        Port:         getEnvAsInt("SERVER_PORT", 8080),
        ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 30),
        WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 30),
        DebugMode:    getEnvAsBool("DEBUG_MODE", false),
    }

    config := &Config{
        Database: dbConfig,
        Server:   serverConfig,
        LogLevel: strings.ToUpper(getEnv("LOG_LEVEL", "INFO")),
    }

    if err := validateConfig(config); err != nil {
        return nil, err
    }

    return config, nil
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
    if value, err := strconv.ParseBool(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func validateConfig(config *Config) error {
    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return &ConfigError{Field: "DB_PORT", Message: "port must be between 1 and 65535"}
    }

    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return &ConfigError{Field: "SERVER_PORT", Message: "port must be between 1 and 65535"}
    }

    if config.Database.Password == "" {
        return &ConfigError{Field: "DB_PASSWORD", Message: "database password cannot be empty"}
    }

    validLogLevels := map[string]bool{"DEBUG": true, "INFO": true, "WARN": true, "ERROR": true}
    if !validLogLevels[config.LogLevel] {
        return &ConfigError{Field: "LOG_LEVEL", Message: "invalid log level"}
    }

    return nil
}

type ConfigError struct {
    Field   string
    Message string
}

func (e *ConfigError) Error() string {
    return "config error: field " + e.Field + " - " + e.Message
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
		Username string `yaml:"username" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASS"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
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
	overrideString(&cfg.Database.Username, "DB_USER")
	overrideString(&cfg.Database.Password, "DB_PASS")
	overrideString(&cfg.Logging.Level, "LOG_LEVEL")
	overrideString(&cfg.Logging.Output, "LOG_OUTPUT")
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
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
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
		return errors.New("server host cannot be empty")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if c.Database.Name == "" {
		return errors.New("database name cannot be empty")
	}
	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}
	if c.Logging.Output == "" {
		c.Logging.Output = "stdout"
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
    DBHost     string
    DBPort     int
    DebugMode  bool
    AllowedIPs []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort: getEnvAsInt("SERVER_PORT", 8080),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnvAsInt("DB_PORT", 5432),
        DebugMode:  getEnvAsBool("DEBUG_MODE", false),
        AllowedIPs: getEnvAsSlice("ALLOWED_IPS", []string{"127.0.0.1"}, ","),
    }

    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return nil, &ConfigError{Field: "SERVER_PORT", Value: strconv.Itoa(cfg.ServerPort)}
    }

    if cfg.DBPort < 1 || cfg.DBPort > 65535 {
        return nil, &ConfigError{Field: "DB_PORT", Value: strconv.Itoa(cfg.DBPort)}
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

type ConfigError struct {
    Field string
    Value string
}

func (e *ConfigError) Error() string {
    return "invalid configuration value for " + e.Field + ": " + e.Value
}
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"SERVER_DEBUG"`
}

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	LogLevel string         `json:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	var config AppConfig

	if configPath != "" {
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			return nil, fmt.Errorf("invalid config path: %w", err)
		}

		data, err := os.ReadFile(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse config JSON: %w", err)
		}
	}

	if err := loadEnvOverrides(&config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadEnvOverrides(config *AppConfig) error {
	overrideStringField(&config.LogLevel, "LOG_LEVEL")

	overrideStringField(&config.Database.Host, "DB_HOST")
	overrideIntField(&config.Database.Port, "DB_PORT")
	overrideStringField(&config.Database.Username, "DB_USER")
	overrideStringField(&config.Database.Password, "DB_PASS")
	overrideStringField(&config.Database.SSLMode, "DB_SSL_MODE")

	overrideIntField(&config.Server.Port, "SERVER_PORT")
	overrideIntField(&config.Server.ReadTimeout, "SERVER_READ_TIMEOUT")
	overrideIntField(&config.Server.WriteTimeout, "SERVER_WRITE_TIMEOUT")
	overrideBoolField(&config.Server.DebugMode, "SERVER_DEBUG")

	return nil
}

func overrideStringField(field *string, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val
	}
}

func overrideIntField(field *int, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		var intVal int
		if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
			*field = intVal
		}
	}
}

func overrideBoolField(field *bool, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		lowerVal := strings.ToLower(val)
		*field = lowerVal == "true" || lowerVal == "1" || lowerVal == "yes"
	}
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}
	if config.Server.ReadTimeout < 0 {
		return fmt.Errorf("server read timeout cannot be negative")
	}
	if config.Server.WriteTimeout < 0 {
		return fmt.Errorf("server write timeout cannot be negative")
	}

	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[strings.ToLower(config.LogLevel)] {
		return fmt.Errorf("invalid log level: %s", config.LogLevel)
	}

	return nil
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
		Driver   string `yaml:"driver" env:"DB_DRIVER"`
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		Username string `yaml:"username" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASS"`
		Name     string `yaml:"name" env:"DB_NAME"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
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
	overrideString(&cfg.Database.Driver, "DB_DRIVER")
	overrideString(&cfg.Database.Host, "DB_HOST")
	overrideInt(&cfg.Database.Port, "DB_PORT")
	overrideString(&cfg.Database.Username, "DB_USER")
	overrideString(&cfg.Database.Password, "DB_PASS")
	overrideString(&cfg.Database.Name, "DB_NAME")
	overrideString(&cfg.Logging.Level, "LOG_LEVEL")
	overrideString(&cfg.Logging.Output, "LOG_OUTPUT")
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