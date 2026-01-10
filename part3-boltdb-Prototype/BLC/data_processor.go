
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	ServerAddress string `json:"server_address"`
	Port          int    `json:"port"`
	EnableLogging bool   `json:"enable_logging"`
	MaxConnections int   `json:"max_connections"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if config.ServerAddress == "" {
		return nil, fmt.Errorf("server_address cannot be empty")
	}
	if config.Port <= 0 || config.Port > 65535 {
		return nil, fmt.Errorf("port must be between 1 and 65535")
	}
	if config.MaxConnections < 1 {
		return nil, fmt.Errorf("max_connections must be at least 1")
	}

	return &config, nil
}

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration loaded successfully:\n")
	fmt.Printf("Server: %s:%d\n", config.ServerAddress, config.Port)
	fmt.Printf("Logging enabled: %v\n", config.EnableLogging)
	fmt.Printf("Max connections: %d\n", config.MaxConnections)
}