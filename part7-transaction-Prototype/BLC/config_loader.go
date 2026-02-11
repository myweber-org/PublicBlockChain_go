package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
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
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func ValidateConfig(config *AppConfig) bool {
	if config.Server.Host == "" || config.Server.Port == 0 {
		log.Println("Invalid server configuration")
		return false
	}
	if config.Database.Host == "" || config.Database.Name == "" {
		log.Println("Invalid database configuration")
		return false
	}
	return true
}
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	Server struct {
		Host string `json:"host" env:"SERVER_HOST" default:"localhost"`
		Port int    `json:"port" env:"SERVER_PORT" default:"8080"`
	} `json:"server"`
	Database struct {
		URL      string `json:"url" env:"DATABASE_URL" required:"true"`
		PoolSize int    `json:"pool_size" env:"DB_POOL_SIZE" default:"10"`
	} `json:"database"`
	LogLevel string `json:"log_level" env:"LOG_LEVEL" default:"info"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	
	if configPath != "" {
		if err := loadFromFile(configPath, cfg); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}
	
	if err := loadFromEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}
	
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return cfg, nil
}

func loadFromFile(path string, cfg *Config) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid config path: %w", err)
	}
	
	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config JSON: %w", err)
	}
	
	return nil
}

func loadFromEnv(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	return processStruct(v, "")
}

func processStruct(v reflect.Value, prefix string) error {
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		if field.Kind() == reflect.Struct {
			newPrefix := prefix
			if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				if parts[0] != "" && parts[0] != "-" {
					newPrefix = strings.TrimPrefix(prefix+"_"+strings.ToUpper(parts[0]), "_")
				}
			}
			if err := processStruct(field, newPrefix); err != nil {
				return err
			}
			continue
		}
		
		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}
		
		envValue := os.Getenv(envTag)
		if envValue == "" {
			defaultTag := fieldType.Tag.Get("default")
			if defaultTag != "" && field.Interface() == reflect.Zero(field.Type()).Interface() {
				envValue = defaultTag
			}
		}
		
		if envValue != "" {
			if err := setFieldValue(field, envValue); err != nil {
				return fmt.Errorf("failed to set field %s: %w", fieldType.Name, err)
			}
		}
	}
	
	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		var intVal int
		if _, err := fmt.Sscanf(value, "%d", &intVal); err != nil {
			return fmt.Errorf("invalid integer value: %s", value)
		}
		field.SetInt(int64(intVal))
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

func validateConfig(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	return validateStruct(v)
}

func validateStruct(v reflect.Value) error {
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		if field.Kind() == reflect.Struct {
			if err := validateStruct(field); err != nil {
				return err
			}
			continue
		}
		
		requiredTag := fieldType.Tag.Get("required")
		if requiredTag == "true" {
			zeroValue := reflect.Zero(field.Type()).Interface()
			if field.Interface() == zeroValue {
				envTag := fieldType.Tag.Get("env")
				return fmt.Errorf("required field %s is not set (environment variable: %s)", fieldType.Name, envTag)
			}
		}
	}
	
	return nil
}package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DebugMode:  getEnvAsBool("DEBUG_MODE", false),
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
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	strValue := getEnv(key, "")
	if value, err := strconv.ParseBool(strValue); err == nil {
		return value
	}
	return defaultValue
}