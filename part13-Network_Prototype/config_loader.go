package config

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
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
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
	overrideString(&cfg.Database.User, "DB_USER")
	overrideString(&cfg.Database.Password, "DB_PASSWORD")
	overrideString(&cfg.Database.SSLMode, "DB_SSL_MODE")
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
	"encoding/json"
	"os"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST" required:"true"`
	Port     int    `json:"port" env:"DB_PORT" default:"5432"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE" default:"disable"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT" default:"8080"`
	ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT" default:"30"`
	WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT" default:"30"`
	DebugMode    bool   `json:"debug_mode" env:"DEBUG_MODE" default:"false"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL" default:"info"`
}

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Features []string       `json:"features" env:"ENABLED_FEATURES"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	config := &AppConfig{}

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

	loadFromEnv(config)
	
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func loadFromEnv(config *AppConfig) {
	loadStruct(&config.Database)
	loadStruct(&config.Server)
	
	if features := os.Getenv("ENABLED_FEATURES"); features != "" {
		config.Features = strings.Split(features, ",")
	}
}

func loadStruct(target interface{}) {
	// Implementation would use reflection to read struct tags
	// and populate values from environment variables
	// This is a simplified placeholder
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return ConfigError{Field: "database.host", Reason: "required field is empty"}
	}
	
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return ConfigError{Field: "server.port", Reason: "port must be between 1 and 65535"}
	}
	
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[config.Server.LogLevel] {
		return ConfigError{Field: "server.log_level", Reason: "invalid log level"}
	}
	
	return nil
}

type ConfigError struct {
	Field  string
	Reason string
}

func (e ConfigError) Error() string {
	return "config error: " + e.Field + " - " + e.Reason
}