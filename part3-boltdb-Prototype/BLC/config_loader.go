package config

import (
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port         int            `yaml:"port"`
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
    Database     DatabaseConfig `yaml:"database"`
}

func LoadConfig(path string) (*ServerConfig, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("config file not found: %s", path)
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %v", err)
    }

    var config ServerConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %v", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %v", err)
    }

    return &config, nil
}

func validateConfig(config *ServerConfig) error {
    if config.Port <= 0 || config.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Port)
    }

    if config.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }

    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }

    if config.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }

    return nil
}package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string
	APIKey      string
	LogLevel    string
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	content = os.ExpandEnv(content)

	lines := strings.Split(content, "\n")
	cfg := &Config{}

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "DATABASE_URL":
			cfg.DatabaseURL = value
		case "API_KEY":
			cfg.APIKey = value
		case "LOG_LEVEL":
			cfg.LogLevel = value
		}
	}

	return cfg, nil
}package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `json:"host" yaml:"host"`
		Port int    `json:"port" yaml:"port"`
	} `json:"server" yaml:"server"`
	Database struct {
		Driver   string `json:"driver" yaml:"driver"`
		Host     string `json:"host" yaml:"host"`
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
	} `json:"database" yaml:"database"`
	LogLevel string `json:"log_level" yaml:"log_level"`
}

func LoadConfig(filePath string) (*Config, error) {
	if filePath == "" {
		return nil, errors.New("config file path cannot be empty")
	}

	ext := filepath.Ext(filePath)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	switch ext {
	case ".json":
		err = json.Unmarshal(fileData, config)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(fileData, config)
	default:
		return nil, errors.New("unsupported config file format")
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return errors.New("server host is required")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if c.Database.Driver == "" {
		return errors.New("database driver is required")
	}
	return nil
}package config

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
		Format string `yaml:"format" env:"LOG_FORMAT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	overrideFromEnv(&config)

	return &config, nil
}

func overrideFromEnv(config *Config) {
	config.Server.Host = getEnvOrDefault("SERVER_HOST", config.Server.Host)
	config.Server.Port = getEnvIntOrDefault("SERVER_PORT", config.Server.Port)

	config.Database.Host = getEnvOrDefault("DB_HOST", config.Database.Host)
	config.Database.Port = getEnvIntOrDefault("DB_PORT", config.Database.Port)
	config.Database.Name = getEnvOrDefault("DB_NAME", config.Database.Name)
	config.Database.User = getEnvOrDefault("DB_USER", config.Database.User)
	config.Database.Password = getEnvOrDefault("DB_PASSWORD", config.Database.Password)
	config.Database.SSLMode = getEnvOrDefault("DB_SSL_MODE", config.Database.SSLMode)

	config.Logging.Level = getEnvOrDefault("LOG_LEVEL", config.Logging.Level)
	config.Logging.Format = getEnvOrDefault("LOG_FORMAT", config.Logging.Format)
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