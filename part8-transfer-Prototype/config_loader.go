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
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort    int
	DatabaseURL   string
	CacheEnabled  bool
	LogLevel      string
	MaxRetries    int
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	
	port, err := getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}
	cfg.ServerPort = port
	
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is required")
	}
	cfg.DatabaseURL = dbURL
	
	cacheEnabled, err := getEnvBool("CACHE_ENABLED", true)
	if err != nil {
		return nil, err
	}
	cfg.CacheEnabled = cacheEnabled
	
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	cfg.LogLevel = strings.ToLower(logLevel)
	
	maxRetries, err := getEnvInt("MAX_RETRIES", 3)
	if err != nil {
		return nil, err
	}
	cfg.MaxRetries = maxRetries
	
	return cfg, nil
}

func getEnvInt(key string, defaultValue int) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue, nil
	}
	
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, errors.New("invalid integer value for " + key)
	}
	return intVal, nil
}

func getEnvBool(key string, defaultValue bool) (bool, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue, nil
	}
	
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return false, errors.New("invalid boolean value for " + key)
	}
	return boolVal, nil
}package config

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
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
	Port         int    `yaml:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	DebugMode    bool   `yaml:"debug_mode" env:"SERVER_DEBUG"`
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Features struct {
		CacheEnabled bool `yaml:"cache_enabled" env:"CACHE_ENABLED"`
		APIRateLimit int  `yaml:"api_rate_limit" env:"API_RATE_LIMIT"`
	} `yaml:"features"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
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

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := overrideFromEnv(&config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) error {
	overrideString := func(field *string, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			*field = val
		}
	}

	overrideInt := func(field *int, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			var intVal int
			if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
				*field = intVal
			}
		}
	}

	overrideBool := func(field *bool, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			*field = val == "true" || val == "1" || val == "yes"
		}
	}

	overrideString(&config.Database.Host, "DB_HOST")
	overrideInt(&config.Database.Port, "DB_PORT")
	overrideString(&config.Database.Username, "DB_USER")
	overrideString(&config.Database.Password, "DB_PASS")
	overrideString(&config.Database.Database, "DB_NAME")
	overrideString(&config.Database.SSLMode, "DB_SSL_MODE")

	overrideInt(&config.Server.Port, "SERVER_PORT")
	overrideInt(&config.Server.ReadTimeout, "SERVER_READ_TIMEOUT")
	overrideInt(&config.Server.WriteTimeout, "SERVER_WRITE_TIMEOUT")
	overrideBool(&config.Server.DebugMode, "SERVER_DEBUG")
	overrideString(&config.Server.LogLevel, "LOG_LEVEL")

	overrideBool(&config.Features.CacheEnabled, "CACHE_ENABLED")
	overrideInt(&config.Features.APIRateLimit, "API_RATE_LIMIT")

	return nil
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("invalid database port")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("invalid server port")
	}
	if config.Server.ReadTimeout < 0 {
		return errors.New("read timeout cannot be negative")
	}
	if config.Server.WriteTimeout < 0 {
		return errors.New("write timeout cannot be negative")
	}
	return nil
}package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
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
    ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
    Debug        bool   `yaml:"debug" env:"SERVER_DEBUG"`
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
    overrideString(&config.Database.Host, "DB_HOST")
    overrideInt(&config.Database.Port, "DB_PORT")
    overrideString(&config.Database.Username, "DB_USER")
    overrideString(&config.Database.Password, "DB_PASS")
    overrideString(&config.Database.Name, "DB_NAME")

    overrideInt(&config.Server.Port, "SERVER_PORT")
    overrideInt(&config.Server.ReadTimeout, "SERVER_READ_TIMEOUT")
    overrideInt(&config.Server.WriteTimeout, "SERVER_WRITE_TIMEOUT")
    overrideBool(&config.Server.Debug, "SERVER_DEBUG")

    overrideString(&config.LogLevel, "LOG_LEVEL")
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
    Debug        bool   `yaml:"debug" env:"DEBUG"`
}

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(configPath string) (*Config, error) {
    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("invalid config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config.Database)
    overrideFromEnv(&config.Server)

    return &config, nil
}

func overrideFromEnv(config interface{}) {
    // Environment variable override logic would be implemented here
    // This is a placeholder for the actual implementation
}package config

import (
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
}

func LoadConfig(filename string) (*Config, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return &config, nil
}

func (c *Config) Validate() error {
    if c.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if c.Server.Port <= 0 {
        return fmt.Errorf("server port must be positive")
    }
    if c.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }
    return nil
}