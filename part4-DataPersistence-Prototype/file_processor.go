package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Server string `json:"server"`
	Port   int    `json:"port"`
	Debug  bool   `json:"debug"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(filename string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func main() {
	config := &Config{
		Server: "api.example.com",
		Port:   8080,
		Debug:  true,
	}

	err := SaveConfig("config.json", config)
	if err != nil {
		log.Fatal(err)
	}

	loadedConfig, err := LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Loaded config: %+v", loadedConfig)

	os.Remove("config.json")
}