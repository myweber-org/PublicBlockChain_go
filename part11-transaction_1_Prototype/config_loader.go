package config

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
    SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
    Port         int    `json:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
    DebugMode    bool   `json:"debug_mode" env:"DEBUG_MODE"`
}

type Config struct {
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    LogLevel string         `json:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var config Config
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return nil, err
    }

    overrideFromEnv(&config)

    if err := validateConfig(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

func overrideFromEnv(config *Config) {
    overrideStruct(config)
}

func overrideStruct(s interface{}) {
    val := reflect.ValueOf(s).Elem()
    typ := val.Type()

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)

        if field.Kind() == reflect.Struct {
            overrideStruct(field.Addr().Interface())
            continue
        }

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
            if intVal, err := strconv.Atoi(envValue); err == nil {
                field.SetInt(int64(intVal))
            }
        case reflect.Bool:
            boolVal := strings.ToLower(envValue) == "true" || envValue == "1"
            field.SetBool(boolVal)
        }
    }
}

func validateConfig(config *Config) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Port < 1 || config.Database.Port > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }
    if config.Server.Port < 1 || config.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if config.LogLevel == "" {
        config.LogLevel = "info"
    }
    
    return nil
}package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int    `env:"SERVER_PORT" default:"8080"`
	LogLevel   string `env:"LOG_LEVEL" default:"info"`
	DBHost     string `env:"DB_HOST" default:"localhost"`
	DBPort     int    `env:"DB_PORT" default:"5432"`
	EnableSSL  bool   `env:"ENABLE_SSL" default:"false"`
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

		envValue := os.Getenv(envKey)
		if envValue == "" {
			envValue = defaultVal
		}

		if err := setField(field, envValue); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %w", structField.Name, err)
		}
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
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
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.ServerPort)
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
		return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
	}

	return nil
}

func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}