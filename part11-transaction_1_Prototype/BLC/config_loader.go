
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
		Driver   string `json:"driver" env:"DB_DRIVER" default:"postgres"`
		Host     string `json:"host" env:"DB_HOST" default:"localhost"`
		Port     int    `json:"port" env:"DB_PORT" default:"5432"`
		Name     string `json:"name" env:"DB_NAME" default:"appdb"`
		User     string `json:"user" env:"DB_USER" default:"postgres"`
		Password string `json:"password" env:"DB_PASSWORD"`
		SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE" default:"disable"`
	} `json:"database"`
	Logging struct {
		Level    string `json:"level" env:"LOG_LEVEL" default:"info"`
		FilePath string `json:"file_path" env:"LOG_FILE_PATH"`
	} `json:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	
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
	
	if err := applyEnvironmentVariables(cfg); err != nil {
		return nil, err
	}
	
	if err := applyDefaults(cfg); err != nil {
		return nil, err
	}
	
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

func applyEnvironmentVariables(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	return processStruct(v, "")
}

func processStruct(v reflect.Value, prefix string) error {
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		if field.Kind() == reflect.Struct {
			tag := fieldType.Tag.Get("json")
			if tag != "" && !strings.HasSuffix(tag, ",omitempty") {
				if err := processStruct(field, prefix+tag+"_"); err != nil {
					return err
				}
			} else {
				if err := processStruct(field, prefix); err != nil {
					return err
				}
			}
			continue
		}
		
		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}
		
		if val, exists := os.LookupEnv(envTag); exists {
			switch field.Kind() {
			case reflect.String:
				field.SetString(val)
			case reflect.Int:
				var intVal int64
				if _, err := fmt.Sscanf(val, "%d", &intVal); err != nil {
					return fmt.Errorf("invalid integer value for %s: %w", envTag, err)
				}
				field.SetInt(intVal)
			default:
				return fmt.Errorf("unsupported field type for %s", envTag)
			}
		}
	}
	
	return nil
}

func applyDefaults(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	return applyDefaultsToStruct(v)
}

func applyDefaultsToStruct(v reflect.Value) error {
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		if field.Kind() == reflect.Struct {
			if err := applyDefaultsToStruct(field); err != nil {
				return err
			}
			continue
		}
		
		defaultTag := fieldType.Tag.Get("default")
		if defaultTag == "" {
			continue
		}
		
		if field.Kind() == reflect.String && field.String() == "" {
			field.SetString(defaultTag)
		} else if field.Kind() == reflect.Int && field.Int() == 0 {
			var intVal int64
			if _, err := fmt.Sscanf(defaultTag, "%d", &intVal); err != nil {
				return fmt.Errorf("invalid default integer value: %w", err)
			}
			field.SetInt(intVal)
		}
	}
	
	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}
	
	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", cfg.Database.Port)
	}
	
	validDrivers := map[string]bool{
		"postgres": true,
		"mysql":    true,
		"sqlite":   true,
	}
	
	if !validDrivers[cfg.Database.Driver] {
		return fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}
	
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	
	if !validLogLevels[cfg.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", cfg.Logging.Level)
	}
	
	return nil
}