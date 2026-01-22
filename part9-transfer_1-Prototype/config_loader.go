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

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Version  string         `yaml:"version"`
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
    if envVal := os.Getenv("DB_HOST"); envVal != "" {
        config.Database.Host = envVal
    }
    if envVal := os.Getenv("DB_PORT"); envVal != "" {
        var port int
        if _, err := fmt.Sscanf(envVal, "%d", &port); err == nil {
            config.Database.Port = port
        }
    }
    if envVal := os.Getenv("DB_USER"); envVal != "" {
        config.Database.Username = envVal
    }
    if envVal := os.Getenv("DB_PASS"); envVal != "" {
        config.Database.Password = envVal
    }
    if envVal := os.Getenv("DB_NAME"); envVal != "" {
        config.Database.Name = envVal
    }
    if envVal := os.Getenv("SERVER_PORT"); envVal != "" {
        var port int
        if _, err := fmt.Sscanf(envVal, "%d", &port); err == nil {
            config.Server.Port = port
        }
    }
    if envVal := os.Getenv("READ_TIMEOUT"); envVal != "" {
        var timeout int
        if _, err := fmt.Sscanf(envVal, "%d", &timeout); err == nil {
            config.Server.ReadTimeout = timeout
        }
    }
    if envVal := os.Getenv("WRITE_TIMEOUT"); envVal != "" {
        var timeout int
        if _, err := fmt.Sscanf(envVal, "%d", &timeout); err == nil {
            config.Server.WriteTimeout = timeout
        }
    }
    if envVal := os.Getenv("DEBUG_MODE"); envVal != "" {
        config.Server.DebugMode = envVal == "true" || envVal == "1"
    }
}

func (c *AppConfig) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if c.Database.Port <= 0 || c.Database.Port > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    return nil
}package config

import (
    "fmt"
    "os"
    "strings"

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
    Debug        bool   `yaml:"debug" env:"SERVER_DEBUG"`
    LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Features []string       `yaml:"features"`
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
    overrideStruct(config, "")
}

func overrideStruct(s interface{}, prefix string) {
    val := reflect.ValueOf(s).Elem()
    typ := val.Type()

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)

        envTag := fieldType.Tag.Get("env")
        yamlTag := fieldType.Tag.Get("yaml")

        if envTag == "" && yamlTag == "" {
            if field.Kind() == reflect.Struct {
                nestedPrefix := prefix
                if yamlTag != "" {
                    nestedPrefix = strings.TrimSuffix(prefix+"_"+strings.ToUpper(yamlTag), "_")
                }
                overrideStruct(field.Addr().Interface(), nestedPrefix)
            }
            continue
        }

        if envTag == "" {
            envTag = strings.ToUpper(yamlTag)
        }

        if prefix != "" {
            envTag = prefix + "_" + envTag
        }

        envValue := os.Getenv(envTag)
        if envValue == "" {
            continue
        }

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
        case reflect.Slice:
            if field.Type().Elem().Kind() == reflect.String {
                items := strings.Split(envValue, ",")
                slice := reflect.MakeSlice(field.Type(), len(items), len(items))
                for j, item := range items {
                    slice.Index(j).SetString(strings.TrimSpace(item))
                }
                field.Set(slice)
            }
        }
    }
}

func ValidateConfig(config *AppConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Port == 0 {
        return fmt.Errorf("database port is required")
    }
    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }
    if config.Server.LogLevel == "" {
        config.Server.LogLevel = "info"
    }
    return nil
}package config

import (
    "fmt"
    "io"
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
    LogLevel string `yaml:"log_level"`
}

func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(c *Config) error {
    if c.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if c.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    if c.LogLevel == "" {
        c.LogLevel = "info"
    }
    return nil
}