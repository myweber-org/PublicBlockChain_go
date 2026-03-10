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
    "strconv"
    "strings"
)

type Config struct {
    DatabaseURL  string
    MaxConnections int
    DebugMode    bool
    AllowedHosts []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
        DebugMode:    getEnvAsBool("DEBUG_MODE", false),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost", "127.0.0.1"}),
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
    valueStr := strings.ToLower(getEnv(key, ""))
    if valueStr == "true" || valueStr == "1" {
        return true
    } else if valueStr == "false" || valueStr == "0" {
        return false
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
	"errors"
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
		Username string `yaml:"username" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASS"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
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

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := overrideFromEnv(&cfg); err != nil {
		return nil, err
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func overrideFromEnv(cfg *Config) error {
	envMap := map[string]interface{}{
		"SERVER_HOST": &cfg.Server.Host,
		"SERVER_PORT": &cfg.Server.Port,
		"DB_HOST":     &cfg.Database.Host,
		"DB_PORT":     &cfg.Database.Port,
		"DB_NAME":     &cfg.Database.Name,
		"DB_USER":     &cfg.Database.Username,
		"DB_PASS":     &cfg.Database.Password,
		"LOG_LEVEL":   &cfg.Logging.Level,
		"LOG_OUTPUT":  &cfg.Logging.Output,
	}

	for envVar, target := range envMap {
		if val := os.Getenv(envVar); val != "" {
			switch t := target.(type) {
			case *string:
				*t = val
			case *int:
				var intVal int
				if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
					*t = intVal
				}
			}
		}
	}
	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.Server.Host == "" {
		return errors.New("server host cannot be empty")
	}
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if cfg.Database.Host == "" {
		return errors.New("database host cannot be empty")
	}
	if cfg.Database.Name == "" {
		return errors.New("database name cannot be empty")
	}
	return nil
}package config

import (
	"os"
	"path/filepath"
	"strings"

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
	cfg := &Config{}

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

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	overrideFromEnv(cfg)

	return cfg, nil
}

func overrideFromEnv(cfg *Config) {
	overrideStruct(cfg, "")
}

func overrideStruct(v interface{}, prefix string) {
	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if field.Kind() == reflect.Struct {
			tag := fieldType.Tag.Get("yaml")
			newPrefix := prefix
			if tag != "" {
				if newPrefix != "" {
					newPrefix = newPrefix + "_"
				}
				newPrefix = newPrefix + strings.ToUpper(tag)
			}
			overrideStruct(field.Addr().Interface(), newPrefix)
			continue
		}

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envVar := prefix
		if envVar != "" {
			envVar = envVar + "_"
		}
		envVar = envVar + envTag

		if envValue := os.Getenv(envVar); envValue != "" {
			switch field.Kind() {
			case reflect.String:
				field.SetString(envValue)
			case reflect.Int:
				if intVal, err := strconv.Atoi(envValue); err == nil {
					field.SetInt(int64(intVal))
				}
			case reflect.Bool:
				if boolVal, err := strconv.ParseBool(envValue); err == nil {
					field.SetBool(boolVal)
				}
			}
		}
	}
}