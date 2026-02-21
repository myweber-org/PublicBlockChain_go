package config

import (
    "fmt"
    "io"
    "os"

    "gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port         int    `yaml:"port"`
    ReadTimeout  int    `yaml:"read_timeout"`
    WriteTimeout int    `yaml:"write_timeout"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Debug    bool           `yaml:"debug"`
}

func LoadConfig(path string) (*AppConfig, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(config *AppConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if config.Server.ReadTimeout < 0 {
        return fmt.Errorf("read timeout cannot be negative")
    }
    if config.Server.WriteTimeout < 0 {
        return fmt.Errorf("write timeout cannot be negative")
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
    EnableDebug bool
    AllowedOrigins []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvOrDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnvOrDefault("DATABASE_URL", "postgres://localhost:5432/app")
    
    debugStr := getEnvOrDefault("ENABLE_DEBUG", "false")
    cfg.EnableDebug = strings.ToLower(debugStr) == "true"
    
    originsStr := getEnvOrDefault("ALLOWED_ORIGINS", "http://localhost:3000")
    cfg.AllowedOrigins = strings.Split(originsStr, ",")
    
    if err := validateConfig(cfg); err != nil {
        return nil, err
    }
    
    return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
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
    Field   string
    Message string
}

func (e *ConfigError) Error() string {
    return "config error: " + e.Field + " - " + e.Message
}package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DebugMode  bool   `json:"debug_mode"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	config.ServerPort = replaceWithEnv(config.ServerPort, "SERVER_PORT")
	config.DBHost = replaceWithEnv(config.DBHost, "DB_HOST")
	config.DBPort = replaceWithEnv(config.DBPort, "DB_PORT")

	return &config, nil
}

func replaceWithEnv(value, envKey string) string {
	if strings.HasPrefix(value, "$") {
		envValue := os.Getenv(envKey)
		if envValue != "" {
			return envValue
		}
	}
	return value
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
    DebugMode    bool
    ReadTimeout  int
    WriteTimeout int
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
}

func LoadConfig() (*Config, error) {
    dbConfig, err := loadDatabaseConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load database config: %w", err)
    }

    serverConfig, err := loadServerConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load server config: %w", err)
    }

    return &Config{
        Database: *dbConfig,
        Server:   *serverConfig,
    }, nil
}

func loadDatabaseConfig() (*DatabaseConfig, error) {
    host := getEnvWithDefault("DB_HOST", "localhost")
    port, err := strconv.Atoi(getEnvWithDefault("DB_PORT", "5432"))
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }

    username := getEnvRequired("DB_USERNAME")
    password := getEnvRequired("DB_PASSWORD")
    database := getEnvWithDefault("DB_NAME", "app_database")

    return &DatabaseConfig{
        Host:     host,
        Port:     port,
        Username: username,
        Password: password,
        Database: database,
    }, nil
}

func loadServerConfig() (*ServerConfig, error) {
    port, err := strconv.Atoi(getEnvWithDefault("SERVER_PORT", "8080"))
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }

    debugMode := strings.ToLower(getEnvWithDefault("DEBUG_MODE", "false")) == "true"
    readTimeout, err := strconv.Atoi(getEnvWithDefault("READ_TIMEOUT", "30"))
    if err != nil {
        return nil, fmt.Errorf("invalid READ_TIMEOUT: %w", err)
    }

    writeTimeout, err := strconv.Atoi(getEnvWithDefault("WRITE_TIMEOUT", "30"))
    if err != nil {
        return nil, fmt.Errorf("invalid WRITE_TIMEOUT: %w", err)
    }

    return &ServerConfig{
        Port:         port,
        DebugMode:    debugMode,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
    }, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvRequired(key string) string {
    value := os.Getenv(key)
    if value == "" {
        panic(fmt.Sprintf("required environment variable %s is not set", key))
    }
    return value
}