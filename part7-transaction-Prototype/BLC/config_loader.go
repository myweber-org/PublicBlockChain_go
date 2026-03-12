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
}package config

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
		Username string `yaml:"username" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASS"`
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
		"SERVER_HOST": &cfg.Server.Host,
		"DB_HOST":     &cfg.Database.Host,
		"DB_NAME":     &cfg.Database.Name,
		"DB_USER":     &cfg.Database.Username,
		"DB_PASS":     &cfg.Database.Password,
		"LOG_LEVEL":   &cfg.Logging.Level,
		"LOG_OUTPUT":  &cfg.Logging.Output,
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
			parsed, err := parseInt(val)
			if err != nil {
				return err
			}
			*field = parsed
		}
	}

	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.Server.Host == "" {
		return errors.New("server host is required")
	}
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
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
    "fmt"
    "io"
    "os"

    "gopkg.in/yaml.v3"
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
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
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
    if c.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    if c.LogLevel == "" {
        c.LogLevel = "info"
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
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    MaxConnections int
    LogLevel string
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{
        ServerPort:     8080,
        DatabaseURL:    "localhost:5432",
        CacheEnabled:   true,
        MaxConnections: 100,
        LogLevel:       "info",
    }

    if portStr := os.Getenv("APP_PORT"); portStr != "" {
        port, err := strconv.Atoi(portStr)
        if err != nil {
            return nil, fmt.Errorf("invalid APP_PORT: %v", err)
        }
        cfg.ServerPort = port
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if cacheStr := os.Getenv("CACHE_ENABLED"); cacheStr != "" {
        cacheEnabled, err := strconv.ParseBool(cacheStr)
        if err != nil {
            return nil, fmt.Errorf("invalid CACHE_ENABLED: %v", err)
        }
        cfg.CacheEnabled = cacheEnabled
    }

    if maxConnStr := os.Getenv("MAX_CONNECTIONS"); maxConnStr != "" {
        maxConn, err := strconv.Atoi(maxConnStr)
        if err != nil {
            return nil, fmt.Errorf("invalid MAX_CONNECTIONS: %v", err)
        }
        cfg.MaxConnections = maxConn
    }

    if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
        validLevels := map[string]bool{
            "debug": true,
            "info":  true,
            "warn":  true,
            "error": true,
        }
        if !validLevels[strings.ToLower(logLevel)] {
            return nil, fmt.Errorf("invalid LOG_LEVEL: %s", logLevel)
        }
        cfg.LogLevel = strings.ToLower(logLevel)
    }

    return cfg, nil
}