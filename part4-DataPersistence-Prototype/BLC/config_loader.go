package config

import (
    "fmt"
    "io/ioutil"
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
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return &config, nil
}

func (c *AppConfig) Validate() error {
    if c.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }
    if c.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    return nil
}package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    DebugMode bool
    AllowedHosts []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnv("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/app")
    
    debugStr := getEnv("DEBUG_MODE", "false")
    cfg.DebugMode = strings.ToLower(debugStr) == "true"
    
    hostsStr := getEnv("ALLOWED_HOSTS", "localhost,127.0.0.1")
    cfg.AllowedHosts = strings.Split(hostsStr, ",")
    
    if err := validateConfig(cfg); err != nil {
        return nil, err
    }
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return &ConfigError{Field: "ServerPort", Message: "port must be between 1 and 65535"}
    }
    
    if cfg.DatabaseURL == "" {
        return &ConfigError{Field: "DatabaseURL", Message: "database URL cannot be empty"}
    }
    
    return nil
}

type ConfigError struct {
    Field string
    Message string
}

func (e *ConfigError) Error() string {
    return "config error: " + e.Field + " - " + e.Message
}package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL  string
	MaxConnections int
	DebugMode    bool
	AllowedHosts []string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL:  getEnv("DB_URL", "postgres://localhost:5432/app"),
		MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
		DebugMode:    getEnvAsBool("DEBUG_MODE", false),
		AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost", "127.0.0.1"}),
	}
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

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
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
	t := reflect.TypeOf(cfg).Elem()
	v := reflect.ValueOf(cfg).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		defaultVal := field.Tag.Get("default")

		if envTag == "" {
			continue
		}

		envValue := os.Getenv(envTag)
		if envValue == "" {
			envValue = defaultVal
		}

		if err := setFieldValue(v.Field(i), envValue); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %w", field.Name, err)
		}
	}

	return cfg, nil
}

func setFieldValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return fmt.Errorf("field cannot be set")
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
		boolVal, err := strconv.ParseBool(strings.ToLower(value))
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}