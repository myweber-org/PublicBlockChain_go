
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int    `env:"SERVER_PORT" default:"8080"`
	DBHost     string `env:"DB_HOST" default:"localhost"`
	DBPort     int    `env:"DB_PORT" default:"5432"`
	DBName     string `env:"DB_NAME" default:"appdb"`
	DebugMode  bool   `env:"DEBUG_MODE" default:"false"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		
		envKey := structField.Tag.Get("env")
		defaultVal := structField.Tag.Get("default")
		
		if envKey == "" {
			continue
		}
		
		envVal := os.Getenv(envKey)
		if envVal == "" {
			envVal = defaultVal
		}
		
		if err := setField(field, envVal); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %w", structField.Name, err)
		}
	}
	
	return cfg, nil
}

func setField(field reflect.Value, value string) error {
	if value == "" {
		return nil
	}
	
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return errors.New("unsupported field type")
	}
	
	return nil
}

func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort    int
    DatabaseURL   string
    LogLevel      string
    CacheEnabled  bool
    MaxWorkers    int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort:    getEnvAsInt("SERVER_PORT", 8080),
        DatabaseURL:   getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        LogLevel:      getEnv("LOG_LEVEL", "info"),
        CacheEnabled:  getEnvAsBool("CACHE_ENABLED", true),
        MaxWorkers:    getEnvAsInt("MAX_WORKERS", 10),
    }

    if err := validateConfig(cfg); err != nil {
        return nil, err
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
    valueStr := strings.ToLower(getEnv(key, ""))
    if valueStr == "true" || valueStr == "1" {
        return true
    } else if valueStr == "false" || valueStr == "0" {
        return false
    }
    return defaultValue
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return &ConfigError{Field: "ServerPort", Message: "port must be between 1 and 65535"}
    }
    if cfg.DatabaseURL == "" {
        return &ConfigError{Field: "DatabaseURL", Message: "database URL cannot be empty"}
    }
    if cfg.MaxWorkers < 1 {
        return &ConfigError{Field: "MaxWorkers", Message: "must have at least 1 worker"}
    }
    return nil
}

type ConfigError struct {
    Field   string
    Message string
}

func (e *ConfigError) Error() string {
    return "config error: " + e.Field + " - " + e.Message
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
    AllowedIPs []string
}

func LoadConfig() (*Config, error) {
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

    cfg.AllowedIPs = getStringSliceEnv("ALLOWED_IPS", []string{"127.0.0.1"})

    if err := validateConfig(cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

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
        intValue, err := strconv.Atoi(value)
        if err != nil {
            return 0, fmt.Errorf("cannot convert %s to integer: %w", key, err)
        }
        return intValue, nil
    }
    return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) (bool, error) {
    if value := os.Getenv(key); value != "" {
        boolValue, err := strconv.ParseBool(value)
        if err != nil {
            return false, fmt.Errorf("cannot convert %s to boolean: %w", key, err)
        }
        return boolValue, nil
    }
    return defaultValue, nil
}

func getStringSliceEnv(key string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        return strings.Split(value, ",")
    }
    return defaultValue
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    
    if cfg.DBPort < 1 || cfg.DBPort > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }
    
    if cfg.DBHost == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    
    return nil
}package config

import (
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(c *Config) error {
    if c.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if c.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }
    return nil
}package config

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
}

type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Debug    bool
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "appdb"),
		},
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 30),
			WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 30),
		},
		Debug: getEnvAsBool("DEBUG", false),
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
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
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func validateConfig(cfg *Config) error {
	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return &ConfigError{Field: "DB_PORT", Message: "port must be between 1 and 65535"}
	}

	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return &ConfigError{Field: "SERVER_PORT", Message: "port must be between 1 and 65535"}
	}

	if strings.TrimSpace(cfg.Database.Password) == "" {
		return &ConfigError{Field: "DB_PASSWORD", Message: "database password cannot be empty"}
	}

	return nil
}

type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return "config error [" + e.Field + "]: " + e.Message
}package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     int    `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASS"`
	Database string `yaml:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `yaml:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
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

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	overrideFromEnv(&config)

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) {
	overrideString(&config.Server.LogLevel, "LOG_LEVEL")
	overrideInt(&config.Server.Port, "SERVER_PORT")
	overrideBool(&config.Server.DebugMode, "DEBUG_MODE")
	overrideString(&config.Database.Host, "DB_HOST")
	overrideInt(&config.Database.Port, "DB_PORT")
	overrideString(&config.Database.Username, "DB_USER")
	overrideString(&config.Database.Password, "DB_PASS")
	overrideString(&config.Database.Database, "DB_NAME")
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
}

func validateConfig(config *AppConfig) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	if config.Database.Host == "" {
		return errors.New("database host is required")
	}

	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if config.Database.Database == "" {
		return errors.New("database name is required")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}

	if !validLogLevels[config.Server.LogLevel] {
		return errors.New("invalid log level")
	}

	return nil
}package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    Port        int
    DatabaseURL string
    DebugMode   bool
    AllowedHosts []string
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{}

    portStr := getEnvWithDefault("APP_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid port value: %v", err)
    }
    cfg.Port = port

    dbURL := getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/appdb")
    if !strings.HasPrefix(dbURL, "postgres://") {
        return nil, fmt.Errorf("invalid database URL format")
    }
    cfg.DatabaseURL = dbURL

    debugStr := getEnvWithDefault("DEBUG_MODE", "false")
    debug, err := strconv.ParseBool(debugStr)
    if err != nil {
        return nil, fmt.Errorf("invalid debug mode value: %v", err)
    }
    cfg.DebugMode = debug

    hostsStr := getEnvWithDefault("ALLOWED_HOSTS", "localhost,127.0.0.1")
    cfg.AllowedHosts = strings.Split(hostsStr, ",")

    if err := validateConfig(cfg); err != nil {
        return nil, err
    }

    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func validateConfig(cfg *AppConfig) error {
    if cfg.Port < 1 || cfg.Port > 65535 {
        return fmt.Errorf("port must be between 1 and 65535")
    }

    if len(cfg.DatabaseURL) == 0 {
        return fmt.Errorf("database URL cannot be empty")
    }

    if len(cfg.AllowedHosts) == 0 {
        return fmt.Errorf("at least one allowed host must be specified")
    }

    return nil
}