package config

import (
	"encoding/json"
	"errors"
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
	Env      string         `json:"env" env:"APP_ENV"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	var config AppConfig

	if configPath != "" {
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			return nil, err
		}

		fileData, err := os.ReadFile(absPath)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(fileData, &config); err != nil {
			return nil, err
		}
	}

	if err := loadFromEnv(&config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadFromEnv(config *AppConfig) error {
	loadStruct := func(s interface{}) error {
		v := reflect.ValueOf(s).Elem()
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := t.Field(i)

			envTag := fieldType.Tag.Get("env")
			if envTag == "" {
				continue
			}

			envValue := os.Getenv(envTag)
			if envValue == "" {
				continue
			}

			switch field.Kind() {
			case reflect.String:
				field.SetString(envValue)
			case reflect.Int:
				intVal, err := strconv.Atoi(envValue)
				if err != nil {
					return err
				}
				field.SetInt(int64(intVal))
			case reflect.Bool:
				boolVal := strings.ToLower(envValue) == "true"
				field.SetBool(boolVal)
			}
		}
		return nil
	}

	if err := loadStruct(&config.Database); err != nil {
		return err
	}
	if err := loadStruct(&config.Server); err != nil {
		return err
	}
	if err := loadStruct(config); err != nil {
		return err
	}

	return nil
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port < 1 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Server.ReadTimeout < 0 {
		return errors.New("server read timeout must be non-negative")
	}
	if config.Server.WriteTimeout < 0 {
		return errors.New("server write timeout must be non-negative")
	}

	return nil
}package config

import (
	"encoding/json"
	"os"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	Database string `json:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}
	
	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, err
		}
	}
	
	overrideFromEnv(config)
	
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	
	return config, nil
}

func overrideFromEnv(config *Config) {
	overrideStruct(&config.Database)
	overrideStruct(&config.Server)
}

func overrideStruct(s interface{}) {
	// Implementation would use reflection to read struct tags
	// and override values from environment variables
	// Simplified for this example
}

func validateConfig(config *Config) error {
	if config.Database.Host == "" {
		return NewValidationError("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return NewValidationError("database port must be between 1 and 65535")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return NewValidationError("server port must be between 1 and 65535")
	}
	if !isValidLogLevel(config.Server.LogLevel) {
		return NewValidationError("invalid log level")
	}
	return nil
}

func isValidLogLevel(level string) bool {
	validLevels := []string{"debug", "info", "warn", "error", "fatal"}
	for _, valid := range validLevels {
		if strings.ToLower(level) == valid {
			return true
		}
	}
	return false
}

type ValidationError struct {
	Message string
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{Message: msg}
}

func (e *ValidationError) Error() string {
	return "config validation error: " + e.Message
}package config

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
	ServerPort int    `json:"server_port" env:"SERVER_PORT"`
	DBHost     string `json:"db_host" env:"DB_HOST"`
	DBPort     int    `json:"db_port" env:"DB_PORT"`
	DebugMode  bool   `json:"debug_mode" env:"DEBUG_MODE"`
	MaxWorkers int    `json:"max_workers" env:"MAX_WORKERS"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}
	
	if configPath != "" {
		if err := loadFromFile(configPath, config); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}
	
	if err := loadFromEnv(config); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}
	
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return config, nil
}

func loadFromFile(path string, config *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	return decoder.Decode(config)
}

func loadFromEnv(config *Config) error {
	configValue := reflect.ValueOf(config).Elem()
	configType := configValue.Type()
	
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}
		
		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}
		
		fieldValue := configValue.Field(i)
		if err := setFieldValue(fieldValue, envValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", field.Name, err)
		}
	}
	
	return nil
}

func setFieldValue(field reflect.Value, value string) error {
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

func validateConfig(config *Config) error {
	var validationErrors []string
	
	if config.ServerPort <= 0 || config.ServerPort > 65535 {
		validationErrors = append(validationErrors, "server_port must be between 1 and 65535")
	}
	
	if config.DBHost == "" {
		validationErrors = append(validationErrors, "db_host is required")
	}
	
	if config.DBPort <= 0 || config.DBPort > 65535 {
		validationErrors = append(validationErrors, "db_port must be between 1 and 65535")
	}
	
	if config.MaxWorkers < 1 {
		validationErrors = append(validationErrors, "max_workers must be at least 1")
	}
	
	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, "; "))
	}
	
	return nil
}