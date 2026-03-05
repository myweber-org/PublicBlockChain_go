package config

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
		return nil, errors.New("config path cannot be empty")
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

	return &config, nil
}

func overrideFromEnv(config *AppConfig) {
	overrideStruct(&config.Server)
	overrideStruct(&config.Database)
}

func overrideStruct(target interface{}) {
	// Environment variable override logic would be implemented here
	// This is a placeholder for the actual implementation
}package config

import (
	"errors"
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

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	var err error

	cfg.ServerPort, err = getIntEnv("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}

	cfg.DBHost = getStringEnv("DB_HOST", "localhost")
	
	cfg.DBPort, err = getIntEnv("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	cfg.DebugMode = getBoolEnv("DEBUG_MODE", false)
	
	apiKeysStr := getStringEnv("API_KEYS", "")
	if apiKeysStr != "" {
		cfg.APIKeys = strings.Split(apiKeysStr, ",")
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
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
			return 0, errors.New("invalid integer value for " + key)
		}
		return intValue, nil
	}
	return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return boolValue
	}
	return defaultValue
}

func validateConfig(cfg *Config) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	
	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	
	if cfg.DBHost == "" {
		return errors.New("database host cannot be empty")
	}
	
	return nil
}
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
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
    LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Version  string         `yaml:"version"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    var config AppConfig

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to get absolute path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    overrideString(&config.Database.Host, "DB_HOST")
    overrideString(&config.Database.Username, "DB_USER")
    overrideString(&config.Database.Password, "DB_PASS")
    overrideString(&config.Database.Name, "DB_NAME")
    overrideInt(&config.Database.Port, "DB_PORT")
    
    overrideInt(&config.Server.Port, "SERVER_PORT")
    overrideInt(&config.Server.ReadTimeout, "READ_TIMEOUT")
    overrideInt(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
    overrideBool(&config.Server.DebugMode, "DEBUG_MODE")
    overrideString(&config.Server.LogLevel, "LOG_LEVEL")
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

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port int    `yaml:"port"`
    Mode string `yaml:"mode"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    if configPath == "" {
        configPath = "config.yaml"
    }

    file, err := os.Open(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    if config.Database.Host == "" {
        config.Database.Host = "localhost"
    }

    if config.Database.Port == 0 {
        config.Database.Port = 5432
    }

    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }

    if config.Server.Mode == "" {
        config.Server.Mode = "development"
    }

    if config.LogLevel == "" {
        config.LogLevel = "info"
    }

    return &config, nil
}

func ValidateConfig(config *AppConfig) error {
    if config.Database.Name == "" {
        return fmt.Errorf("database name is required")
    }

    if config.Database.Username == "" {
        return fmt.Errorf("database username is required")
    }

    if config.Server.Port < 1 || config.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }

    validModes := map[string]bool{
        "development": true,
        "staging":     true,
        "production":  true,
    }

    if !validModes[config.Server.Mode] {
        return fmt.Errorf("invalid server mode: %s", config.Server.Mode)
    }

    validLogLevels := map[string]bool{
        "debug":   true,
        "info":    true,
        "warning": true,
        "error":   true,
        "fatal":   true,
    }

    if !validLogLevels[config.LogLevel] {
        return fmt.Errorf("invalid log level: %s", config.LogLevel)
    }

    return nil
}package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
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
	cfg := &Config{}

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

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	overrideFromEnv(cfg)

	return cfg, nil
}

func overrideFromEnv(cfg *Config) {
	overrideStruct(cfg, "")
}

func overrideStruct(s interface{}, prefix string) {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			newPrefix := prefix
			if tag := fieldType.Tag.Get("yaml"); tag != "" {
				if newPrefix != "" {
					newPrefix = newPrefix + "_"
				}
				newPrefix = newPrefix + strings.ToUpper(tag)
			}
			overrideStruct(field.Addr().Interface(), newPrefix)
			continue
		}

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		if prefix != "" {
			envTag = prefix + "_" + envTag
		}

		if envValue := os.Getenv(envTag); envValue != "" {
			switch field.Kind() {
			case reflect.String:
				field.SetString(envValue)
			case reflect.Int:
				if intVal, err := strconv.Atoi(envValue); err == nil {
					field.SetInt(int64(intVal))
				}
			case reflect.Bool:
				if boolVal, err := strconv.ParseBool(envValue); err == nil {
					field.SetBool(boolVal)
				}
			}
		}
	}
}