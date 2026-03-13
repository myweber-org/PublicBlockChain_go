package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST"`
		Port int    `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		URL      string `yaml:"url" env:"DB_URL"`
		MaxConns int    `yaml:"max_connections" env:"DB_MAX_CONNS"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	overrideWithEnvVars(cfg)
	return cfg, nil
}

func overrideWithEnvVars(cfg *Config) {
	v := reflect.ValueOf(cfg).Elem()
	overrideStruct(v)
}

func overrideStruct(v reflect.Value) {
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		if field.Kind() == reflect.Struct {
			overrideStruct(field)
			continue
		}
		
		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}
		
		if envValue := os.Getenv(envTag); envValue != "" {
			setFieldValue(field, envValue)
		}
	}
}

func setFieldValue(field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		if intVal, err := strconv.Atoi(value); err == nil {
			field.SetInt(int64(intVal))
		}
	case reflect.Bool:
		if boolVal, err := strconv.ParseBool(value); err == nil {
			field.SetBool(boolVal)
		}
	}
}