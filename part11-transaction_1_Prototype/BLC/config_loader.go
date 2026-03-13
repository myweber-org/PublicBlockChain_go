
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	Server struct {
		Host string `json:"host" env:"SERVER_HOST" default:"localhost"`
		Port int    `json:"port" env:"SERVER_PORT" default:"8080"`
	} `json:"server"`
	Database struct {
		Driver   string `json:"driver" env:"DB_DRIVER" default:"postgres"`
		Host     string `json:"host" env:"DB_HOST" default:"localhost"`
		Port     int    `json:"port" env:"DB_PORT" default:"5432"`
		Name     string `json:"name" env:"DB_NAME" default:"appdb"`
		User     string `json:"user" env:"DB_USER" default:"postgres"`
		Password string `json:"password" env:"DB_PASSWORD"`
		SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE" default:"disable"`
	} `json:"database"`
	Logging struct {
		Level    string `json:"level" env:"LOG_LEVEL" default:"info"`
		FilePath string `json:"file_path" env:"LOG_FILE_PATH"`
	} `json:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	
	if configPath != "" {
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			return nil, fmt.Errorf("invalid config path: %w", err)
		}
		
		data, err := os.ReadFile(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config JSON: %w", err)
		}
	}
	
	if err := applyEnvironmentVariables(cfg); err != nil {
		return nil, err
	}
	
	if err := applyDefaults(cfg); err != nil {
		return nil, err
	}
	
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

func applyEnvironmentVariables(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	return processStruct(v, "")
}

func processStruct(v reflect.Value, prefix string) error {
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		if field.Kind() == reflect.Struct {
			tag := fieldType.Tag.Get("json")
			if tag != "" && !strings.HasSuffix(tag, ",omitempty") {
				if err := processStruct(field, prefix+tag+"_"); err != nil {
					return err
				}
			} else {
				if err := processStruct(field, prefix); err != nil {
					return err
				}
			}
			continue
		}
		
		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}
		
		if val, exists := os.LookupEnv(envTag); exists {
			switch field.Kind() {
			case reflect.String:
				field.SetString(val)
			case reflect.Int:
				var intVal int64
				if _, err := fmt.Sscanf(val, "%d", &intVal); err != nil {
					return fmt.Errorf("invalid integer value for %s: %w", envTag, err)
				}
				field.SetInt(intVal)
			default:
				return fmt.Errorf("unsupported field type for %s", envTag)
			}
		}
	}
	
	return nil
}

func applyDefaults(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	return applyDefaultsToStruct(v)
}

func applyDefaultsToStruct(v reflect.Value) error {
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		if field.Kind() == reflect.Struct {
			if err := applyDefaultsToStruct(field); err != nil {
				return err
			}
			continue
		}
		
		defaultTag := fieldType.Tag.Get("default")
		if defaultTag == "" {
			continue
		}
		
		if field.Kind() == reflect.String && field.String() == "" {
			field.SetString(defaultTag)
		} else if field.Kind() == reflect.Int && field.Int() == 0 {
			var intVal int64
			if _, err := fmt.Sscanf(defaultTag, "%d", &intVal); err != nil {
				return fmt.Errorf("invalid default integer value: %w", err)
			}
			field.SetInt(intVal)
		}
	}
	
	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}
	
	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", cfg.Database.Port)
	}
	
	validDrivers := map[string]bool{
		"postgres": true,
		"mysql":    true,
		"sqlite":   true,
	}
	
	if !validDrivers[cfg.Database.Driver] {
		return fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}
	
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	
	if !validLogLevels[cfg.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", cfg.Logging.Level)
	}
	
	return nil
}package config

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

    return nil
}package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	Database string `json:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	var config AppConfig

	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			return nil, fmt.Errorf("failed to decode config: %w", err)
		}
	}

	overrideFromEnv(&config)

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) {
	overrideStruct(&config.Server)
	overrideStruct(&config.Database)
}

func overrideStruct(s interface{}) {
	val := reflect.ValueOf(s).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envValue)
		case reflect.Int:
			var intVal int
			fmt.Sscanf(envValue, "%d", &intVal)
			field.SetInt(int64(intVal))
		case reflect.Bool:
			boolVal := strings.ToLower(envValue) == "true" || envValue == "1"
			field.SetBool(boolVal)
		}
	}
}

func validateConfig(config *AppConfig) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", config.Database.Port)
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[strings.ToLower(config.Server.LogLevel)] {
		return fmt.Errorf("invalid log level: %s", config.Server.LogLevel)
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
        return nil, fmt.Errorf("failed to resolve config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
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