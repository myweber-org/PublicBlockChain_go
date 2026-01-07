package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	APIKeys    []string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	var err error

	cfg.ServerPort, err = getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}

	cfg.DBHost = getEnvString("DB_HOST", "localhost")
	
	cfg.DBPort, err = getEnvInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	cfg.DebugMode, err = getEnvBool("DEBUG_MODE", false)
	if err != nil {
		return nil, err
	}

	apiKeysStr := getEnvString("API_KEYS", "")
	if apiKeysStr != "" {
		cfg.APIKeys = strings.Split(apiKeysStr, ",")
		for i, key := range cfg.APIKeys {
			cfg.APIKeys[i] = strings.TrimSpace(key)
		}
	}

	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return nil, errors.New("invalid server port range")
	}

	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return nil, errors.New("invalid database port range")
	}

	return cfg, nil
}

func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("invalid integer value for " + key)
		}
		return intValue, nil
	}
	return defaultValue, nil
}

func getEnvBool(key string, defaultValue bool) (bool, error) {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return false, errors.New("invalid boolean value for " + key)
		}
		return boolValue, nil
	}
	return defaultValue, nil
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

func LoadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
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
}package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DBHost     string
    DBPort     int
    DebugMode  bool
    MaxWorkers int
}

func Load() (*Config, error) {
    cfg := &Config{
        ServerPort: 8080,
        DBHost:     "localhost",
        DBPort:     5432,
        DebugMode:  false,
        MaxWorkers: 10,
    }

    if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
        port, err := strconv.Atoi(portStr)
        if err != nil {
            return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
        }
        cfg.ServerPort = port
    }

    if host := os.Getenv("DB_HOST"); host != "" {
        cfg.DBHost = host
    }

    if portStr := os.Getenv("DB_PORT"); portStr != "" {
        port, err := strconv.Atoi(portStr)
        if err != nil {
            return nil, fmt.Errorf("invalid DB_PORT: %v", err)
        }
        cfg.DBPort = port
    }

    if debugStr := os.Getenv("DEBUG_MODE"); debugStr != "" {
        debug, err := strconv.ParseBool(strings.ToLower(debugStr))
        if err != nil {
            return nil, fmt.Errorf("invalid DEBUG_MODE: %v", err)
        }
        cfg.DebugMode = debug
    }

    if workersStr := os.Getenv("MAX_WORKERS"); workersStr != "" {
        workers, err := strconv.Atoi(workersStr)
        if err != nil {
            return nil, fmt.Errorf("invalid MAX_WORKERS: %v", err)
        }
        if workers < 1 {
            return nil, fmt.Errorf("MAX_WORKERS must be positive")
        }
        cfg.MaxWorkers = workers
    }

    return cfg, nil
}