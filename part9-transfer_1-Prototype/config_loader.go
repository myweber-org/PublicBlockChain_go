package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    DatabaseURL string
    MaxConnections int
    DebugMode bool
    AllowedOrigins []string
}

func Load() (*Config, error) {
    cfg := &Config{
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
        DebugMode: getEnvAsBool("DEBUG_MODE", false),
        AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"localhost:3000"}, ","),
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
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DebugMode  bool   `json:"debug_mode"`
}

var (
	instance *Config
	once     sync.Once
)

func Load() *Config {
	once.Do(func() {
		instance = &Config{
			ServerPort: getEnv("SERVER_PORT", "8080"),
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DebugMode:  getEnv("DEBUG", "false") == "true",
		}

		configFile := getEnv("CONFIG_FILE", "")
		if configFile != "" {
			loadFromFile(configFile)
		}
	})
	return instance
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadFromFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	var fileConfig Config
	if err := json.Unmarshal(data, &fileConfig); err != nil {
		return
	}

	if fileConfig.ServerPort != "" {
		instance.ServerPort = fileConfig.ServerPort
	}
	if fileConfig.DBHost != "" {
		instance.DBHost = fileConfig.DBHost
	}
	if fileConfig.DBPort != "" {
		instance.DBPort = fileConfig.DBPort
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
	Username string `json:"username" env:"DB_USER" required:"true"`
	Password string `json:"password" env:"DB_PASS" required:"true"`
	Database string `json:"database" env:"DB_NAME" required:"true"`
	SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE" default:"disable"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT" default:"8080"`
	ReadTimeout  int    `json:"read_timeout" env:"SERVER_READ_TIMEOUT" default:"30"`
	WriteTimeout int    `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT" default:"30"`
	DebugMode    bool   `json:"debug_mode" env:"SERVER_DEBUG" default:"false"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL" default:"info"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
}

func LoadConfig(configPath string) (*Config, error) {
	var cfg Config

	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&cfg); err != nil {
			return nil, err
		}
	}

	if err := loadFromEnv(&cfg.Database); err != nil {
		return nil, err
	}
	if err := loadFromEnv(&cfg.Server); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadFromEnv(config interface{}) error {
	// Implementation would use reflection to read struct tags
	// and populate fields from environment variables
	// This is a simplified placeholder
	return nil
}

func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return errors.New("database host is required")
	}
	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	return nil
}