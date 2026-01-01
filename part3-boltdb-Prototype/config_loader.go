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
}