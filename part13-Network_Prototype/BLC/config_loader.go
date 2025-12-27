package config

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
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
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
    if val := os.Getenv("DB_HOST"); val != "" {
        config.Database.Host = val
    }
    if val := os.Getenv("DB_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Database.Port)
    }
    if val := os.Getenv("DB_USER"); val != "" {
        config.Database.Username = val
    }
    if val := os.Getenv("DB_PASS"); val != "" {
        config.Database.Password = val
    }
    if val := os.Getenv("DB_NAME"); val != "" {
        config.Database.Name = val
    }
    if val := os.Getenv("SERVER_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Server.Port)
    }
    if val := os.Getenv("READ_TIMEOUT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Server.ReadTimeout)
    }
    if val := os.Getenv("WRITE_TIMEOUT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Server.WriteTimeout)
    }
    if val := os.Getenv("DEBUG_MODE"); val != "" {
        config.Server.DebugMode = val == "true" || val == "1"
    }
    if val := os.Getenv("LOG_LEVEL"); val != "" {
        config.LogLevel = val
    }
}
package config

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
    SSLMode  string
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
    DebugMode    bool
}

type AppConfig struct {
    Database DatabaseConfig
    Server   ServerConfig
    LogLevel string
}

func LoadConfig() (*AppConfig, error) {
    dbConfig, err := loadDatabaseConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load database config: %w", err)
    }

    serverConfig, err := loadServerConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load server config: %w", err)
    }

    logLevel := getEnvWithDefault("LOG_LEVEL", "info")

    return &AppConfig{
        Database: *dbConfig,
        Server:   *serverConfig,
        LogLevel: logLevel,
    }, nil
}

func loadDatabaseConfig() (*DatabaseConfig, error) {
    host := getEnvRequired("DB_HOST")
    portStr := getEnvRequired("DB_PORT")
    username := getEnvRequired("DB_USERNAME")
    password := getEnvRequired("DB_PASSWORD")
    database := getEnvRequired("DB_NAME")
    sslMode := getEnvWithDefault("DB_SSL_MODE", "require")

    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT value: %s", portStr)
    }

    if port < 1 || port > 65535 {
        return nil, fmt.Errorf("DB_PORT must be between 1 and 65535")
    }

    return &DatabaseConfig{
        Host:     host,
        Port:     port,
        Username: username,
        Password: password,
        Database: database,
        SSLMode:  sslMode,
    }, nil
}

func loadServerConfig() (*ServerConfig, error) {
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    readTimeoutStr := getEnvWithDefault("READ_TIMEOUT", "30")
    writeTimeoutStr := getEnvWithDefault("WRITE_TIMEOUT", "30")
    debugModeStr := getEnvWithDefault("DEBUG_MODE", "false")

    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT value: %s", portStr)
    }

    readTimeout, err := strconv.Atoi(readTimeoutStr)
    if err != nil {
        return nil, fmt.Errorf("invalid READ_TIMEOUT value: %s", readTimeoutStr)
    }

    writeTimeout, err := strconv.Atoi(writeTimeoutStr)
    if err != nil {
        return nil, fmt.Errorf("invalid WRITE_TIMEOUT value: %s", writeTimeoutStr)
    }

    debugMode := strings.ToLower(debugModeStr) == "true"

    if port < 1 || port > 65535 {
        return nil, fmt.Errorf("SERVER_PORT must be between 1 and 65535")
    }

    return &ServerConfig{
        Port:         port,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
        DebugMode:    debugMode,
    }, nil
}

func getEnvRequired(key string) string {
    value := os.Getenv(key)
    if value == "" {
        panic(fmt.Sprintf("required environment variable %s is not set", key))
    }
    return value
}

func getEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}