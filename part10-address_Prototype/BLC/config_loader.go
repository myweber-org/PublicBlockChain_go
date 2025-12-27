package config

import (
    "fmt"
    "io/ioutil"
    "os"

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
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("config file does not exist: %s", path)
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %v", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %v", err)
    }

    if config.Server.Host == "" {
        return nil, fmt.Errorf("server host cannot be empty")
    }
    if config.Server.Port <= 0 {
        return nil, fmt.Errorf("server port must be positive")
    }
    if config.Database.Name == "" {
        return nil, fmt.Errorf("database name cannot be empty")
    }

    return &config, nil
}package config

import (
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Name     string `yaml:"name"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
    } `yaml:"database"`
    LogLevel string `yaml:"log_level"`
}

func LoadConfig(path string) (*Config, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("config file not found: %s", path)
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %v", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %v", err)
    }

    if err := validateConfig(&cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %v", err)
    }

    return &cfg, nil
}

func validateConfig(cfg *Config) error {
    if cfg.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
    }
    if cfg.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    if cfg.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }
    if cfg.LogLevel == "" {
        cfg.LogLevel = "info"
    }

    return nil
}
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
		Username string `json:"username" env:"DB_USERNAME"`
		Password string `json:"password" env:"DB_PASSWORD"`
		SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE" default:"disable"`
	} `json:"database"`
	Logging struct {
		Level  string `json:"level" env:"LOG_LEVEL" default:"info"`
		Format string `json:"format" env:"LOG_FORMAT" default:"json"`
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
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}
	
	if err := loadFromEnv(cfg); err != nil {
		return nil, err
	}
	
	if err := setDefaults(cfg); err != nil {
		return nil, err
	}
	
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

func loadFromEnv(cfg *Config) error {
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
			if tag != "" {
				tag = strings.Split(tag, ",")[0]
			}
			newPrefix := prefix
			if tag != "" {
				if newPrefix != "" {
					newPrefix = newPrefix + "_"
				}
				newPrefix = newPrefix + strings.ToUpper(tag)
			}
			if err := processStruct(field, newPrefix); err != nil {
				return err
			}
			continue
		}
		
		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}
		
		if prefix != "" {
			envTag = prefix + "_" + envTag
		}
		
		if value, exists := os.LookupEnv(envTag); exists {
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Int:
				var intVal int64
				if _, err := fmt.Sscanf(value, "%d", &intVal); err != nil {
					return fmt.Errorf("invalid integer value for %s: %s", envTag, value)
				}
				field.SetInt(intVal)
			default:
				return fmt.Errorf("unsupported field type for %s: %v", envTag, field.Kind())
			}
		}
	}
	
	return nil
}

func setDefaults(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	return setDefaultsRecursive(v)
}

func setDefaultsRecursive(v reflect.Value) error {
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		if field.Kind() == reflect.Struct {
			if err := setDefaultsRecursive(field); err != nil {
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
				return fmt.Errorf("invalid default integer value for %s: %s", fieldType.Name, defaultTag)
			}
			field.SetInt(intVal)
		}
	}
	
	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.Server.Host == "" {
		return fmt.Errorf("server host cannot be empty")
	}
	
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}
	
	if cfg.Database.Driver == "" {
		return fmt.Errorf("database driver cannot be empty")
	}
	
	if cfg.Database.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}
	
	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}
	
	if cfg.Database.Name == "" {
		return fmt.Errorf("database name cannot be empty")
	}
	
	if cfg.Database.Username == "" {
		return fmt.Errorf("database username cannot be empty")
	}
	
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	
	if !validLogLevels[strings.ToLower(cfg.Logging.Level)] {
		return fmt.Errorf("invalid log level: %s", cfg.Logging.Level)
	}
	
	validLogFormats := map[string]bool{
		"json": true,
		"text": true,
	}
	
	if !validLogFormats[strings.ToLower(cfg.Logging.Format)] {
		return fmt.Errorf("invalid log format: %s", cfg.Logging.Format)
	}
	
	return nil
}

func (c *Config) GetDSN() string {
	db := c.Database
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		db.Driver,
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.SSLMode,
	)
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}