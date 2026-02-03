package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	ServerPort string `json:"server_port" env:"SERVER_PORT"`
	DBHost     string `json:"db_host" env:"DB_HOST"`
	DBPort     int    `json:"db_port" env:"DB_PORT"`
	DebugMode  bool   `json:"debug_mode" env:"DEBUG_MODE"`
}

func LoadConfig(configPath string) (*Config, error) {
	var cfg Config

	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	}

	cfg.loadFromEnv()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) loadFromEnv() {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		c.ServerPort = port
	}
	if host := os.Getenv("DB_HOST"); host != "" {
		c.DBHost = host
	}
}

func (c *Config) validate() error {
	if c.ServerPort == "" {
		return errors.New("server port is required")
	}
	if c.DBHost == "" {
		return errors.New("database host is required")
	}
	if c.DBPort <= 0 || c.DBPort > 65535 {
		return errors.New("invalid database port")
	}
	return nil
}

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}package config

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
    DebugMode    bool   `yaml:"debug" env:"DEBUG_MODE"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)
    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    setFieldFromEnv(&config.Database.Host, "DB_HOST")
    setFieldFromEnv(&config.Database.Port, "DB_PORT")
    setFieldFromEnv(&config.Database.Username, "DB_USER")
    setFieldFromEnv(&config.Database.Password, "DB_PASS")
    setFieldFromEnv(&config.Database.Name, "DB_NAME")
    
    setFieldFromEnv(&config.Server.Port, "SERVER_PORT")
    setFieldFromEnv(&config.Server.ReadTimeout, "READ_TIMEOUT")
    setFieldFromEnv(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
    setFieldFromEnv(&config.Server.DebugMode, "DEBUG_MODE")
    
    setFieldFromEnv(&config.LogLevel, "LOG_LEVEL")
}

func setFieldFromEnv(field interface{}, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        switch f := field.(type) {
        case *string:
            *f = val
        case *int:
            fmt.Sscanf(val, "%d", f)
        case *bool:
            *f = val == "true" || val == "1"
        }
    }
}

func DefaultConfigPath() string {
    paths := []string{
        "./config.yaml",
        "./config/config.yaml",
        "/etc/app/config.yaml",
    }
    
    for _, path := range paths {
        if _, err := os.Stat(path); err == nil {
            absPath, _ := filepath.Abs(path)
            return absPath
        }
    }
    
    return ""
}package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DebugMode  bool
    DatabaseURL string
    AllowedHosts []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
        DebugMode:   getEnvAsBool("DEBUG_MODE", false),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost", "127.0.0.1"}, ","),
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
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnv(key, "")
    if val, err := strconv.ParseBool(valueStr); err == nil {
        return val
    }
    return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
    valueStr := getEnv(key, "")
    if valueStr == "" {
        return defaultValue
    }
    return strings.Split(valueStr, sep)
}package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    EnableDebug bool
    MaxConnections int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT value: %v", err)
    }
    cfg.ServerPort = port
    
    dbURL := getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/app")
    if !strings.HasPrefix(dbURL, "postgres://") {
        return nil, fmt.Errorf("invalid DATABASE_URL format")
    }
    cfg.DatabaseURL = dbURL
    
    debugStr := getEnvWithDefault("ENABLE_DEBUG", "false")
    debug, err := strconv.ParseBool(debugStr)
    if err != nil {
        return nil, fmt.Errorf("invalid ENABLE_DEBUG value: %v", err)
    }
    cfg.EnableDebug = debug
    
    maxConnStr := getEnvWithDefault("MAX_CONNECTIONS", "100")
    maxConn, err := strconv.Atoi(maxConnStr)
    if err != nil {
        return nil, fmt.Errorf("invalid MAX_CONNECTIONS value: %v", err)
    }
    if maxConn <= 0 {
        return nil, fmt.Errorf("MAX_CONNECTIONS must be positive")
    }
    cfg.MaxConnections = maxConn
    
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}