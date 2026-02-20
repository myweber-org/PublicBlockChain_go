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
}

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *Config) {
    config.Database.Host = getEnvOrDefault("DB_HOST", config.Database.Host)
    config.Database.Port = getEnvIntOrDefault("DB_PORT", config.Database.Port)
    config.Database.Username = getEnvOrDefault("DB_USER", config.Database.Username)
    config.Database.Password = getEnvOrDefault("DB_PASS", config.Database.Password)
    config.Database.Name = getEnvOrDefault("DB_NAME", config.Database.Name)

    config.Server.Port = getEnvIntOrDefault("SERVER_PORT", config.Server.Port)
    config.Server.ReadTimeout = getEnvIntOrDefault("READ_TIMEOUT", config.Server.ReadTimeout)
    config.Server.WriteTimeout = getEnvIntOrDefault("WRITE_TIMEOUT", config.Server.WriteTimeout)
    config.Server.DebugMode = getEnvBoolOrDefault("DEBUG_MODE", config.Server.DebugMode)

    config.LogLevel = getEnvOrDefault("LOG_LEVEL", config.LogLevel)
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        var result int
        if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
            return result
        }
    }
    return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        return value == "true" || value == "1" || value == "yes"
    }
    return defaultValue
}

func DefaultConfigPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "./config.yaml"
    }
    return filepath.Join(homeDir, ".app", "config.yaml")
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

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    var config AppConfig

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("invalid config path: %w", err)
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
    overrideInt(&config.Database.Port, "DB_PORT")
    overrideString(&config.Database.Username, "DB_USER")
    overrideString(&config.Database.Password, "DB_PASS")
    overrideString(&config.Database.Name, "DB_NAME")
    
    overrideInt(&config.Server.Port, "SERVER_PORT")
    overrideInt(&config.Server.ReadTimeout, "READ_TIMEOUT")
    overrideInt(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
    overrideBool(&config.Server.Debug, "DEBUG")
    
    overrideString(&config.LogLevel, "LOG_LEVEL")
}

func overrideString(field *string, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        *field = val
    }
}

func overrideInt(field *int, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        var temp int
        if _, err := fmt.Sscanf(val, "%d", &temp); err == nil {
            *field = temp
        }
    }
}

func overrideBool(field *bool, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        *field = val == "true" || val == "1" || val == "TRUE"
    }
}package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    EnableLogging bool
    AllowedOrigins []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/app")
    
    loggingStr := getEnvWithDefault("ENABLE_LOGGING", "true")
    cfg.EnableLogging = strings.ToLower(loggingStr) == "true"
    
    originsStr := getEnvWithDefault("ALLOWED_ORIGINS", "http://localhost:3000")
    cfg.AllowedOrigins = strings.Split(originsStr, ",")
    
    if err := validateConfig(cfg); err != nil {
        return nil, err
    }
    
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return strconv.ErrRange
    }
    
    if cfg.DatabaseURL == "" {
        return strconv.ErrSyntax
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
	DebugMode  bool
	DatabaseURL string
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
	
	debugStr := getEnv("DEBUG_MODE", "false")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"
	
	cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/appdb")
	
	hostsStr := getEnv("ALLOWED_HOSTS", "localhost,127.0.0.1")
	cfg.AllowedHosts = strings.Split(hostsStr, ",")
	
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, errors.New("config path cannot be empty")
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	overrideWithEnvVars(&config)

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func overrideWithEnvVars(config *Config) {
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := parseInt(port); err == nil {
			config.Server.Port = p
		}
	}
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := parseInt(port); err == nil {
			config.Database.Port = p
		}
	}
	if user := os.Getenv("DB_USERNAME"); user != "" {
		config.Database.Username = user
	}
	if pass := os.Getenv("DB_PASSWORD"); pass != "" {
		config.Database.Password = pass
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		config.Database.Name = name
	}
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Logging.Level = strings.ToUpper(level)
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		config.Logging.Format = format
	}
}

func validateConfig(config *Config) error {
	if config.Server.Host == "" {
		return errors.New("server host is required")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Database.Name == "" {
		return errors.New("database name is required")
	}
	if config.Logging.Level != "" {
		validLevels := map[string]bool{
			"DEBUG": true,
			"INFO":  true,
			"WARN":  true,
			"ERROR": true,
			"FATAL": true,
		}
		if !validLevels[config.Logging.Level] {
			return errors.New("invalid logging level")
		}
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
    "os"
    "strconv"
    "strings"
)

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
    DebugMode    bool
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    LogLevel string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}

    dbHost := getEnvWithDefault("DB_HOST", "localhost")
    dbPort := getEnvAsInt("DB_PORT", 5432)
    dbUser := getEnvWithDefault("DB_USER", "postgres")
    dbPass := getEnvWithDefault("DB_PASS", "")
    dbName := getEnvWithDefault("DB_NAME", "appdb")

    cfg.Database = DatabaseConfig{
        Host:     dbHost,
        Port:     dbPort,
        Username: dbUser,
        Password: dbPass,
        Database: dbName,
    }

    serverPort := getEnvAsInt("SERVER_PORT", 8080)
    readTimeout := getEnvAsInt("READ_TIMEOUT", 30)
    writeTimeout := getEnvAsInt("WRITE_TIMEOUT", 30)
    debugMode := getEnvAsBool("DEBUG_MODE", false)

    cfg.Server = ServerConfig{
        Port:         serverPort,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
        DebugMode:    debugMode,
    }

    logLevel := getEnvWithDefault("LOG_LEVEL", "info")
    allowedLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !allowedLevels[strings.ToLower(logLevel)] {
        return nil, fmt.Errorf("invalid log level: %s", logLevel)
    }
    cfg.LogLevel = strings.ToLower(logLevel)

    if cfg.Database.Password == "" {
        return nil, fmt.Errorf("database password must be set")
    }

    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return nil, fmt.Errorf("server port must be between 1 and 65535")
    }

    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnvWithDefault(key, "")
    if valueStr == "" {
        return defaultValue
    }
    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnvWithDefault(key, "")
    if valueStr == "" {
        return defaultValue
    }
    value, err := strconv.ParseBool(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}package config

import (
    "encoding/json"
    "errors"
    "os"
    "path/filepath"
)

type Config struct {
    ServerPort int    `json:"server_port"`
    LogLevel   string `json:"log_level"`
    CacheSize  int    `json:"cache_size"`
    EnableTLS  bool   `json:"enable_tls"`
}

func LoadConfig(configPath string) (*Config, error) {
    if configPath == "" {
        configPath = filepath.Join(".", "config.json")
    }

    fileData, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := json.Unmarshal(fileData, &cfg); err != nil {
        return nil, err
    }

    if err := validateConfig(&cfg); err != nil {
        return nil, err
    }

    setDefaults(&cfg)
    return &cfg, nil
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return errors.New("invalid server port")
    }

    validLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }

    if !validLogLevels[cfg.LogLevel] {
        return errors.New("invalid log level")
    }

    if cfg.CacheSize < 0 {
        return errors.New("cache size cannot be negative")
    }

    return nil
}

func setDefaults(cfg *Config) {
    if cfg.ServerPort == 0 {
        cfg.ServerPort = 8080
    }

    if cfg.LogLevel == "" {
        cfg.LogLevel = "info"
    }

    if cfg.CacheSize == 0 {
        cfg.CacheSize = 100
    }
}