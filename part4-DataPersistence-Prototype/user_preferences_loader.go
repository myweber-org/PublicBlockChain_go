package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type UserPreferences struct {
	Theme     string `json:"theme"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
	NotificationsEnabled bool `json:"notifications_enabled"`
}

func LoadPreferences(filename string) (*UserPreferences, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open preferences file: %w", err)
	}
	defer file.Close()

	var prefs UserPreferences
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&prefs); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &prefs, nil
}

func main() {
	prefs, err := LoadPreferences("preferences.json")
	if err != nil {
		fmt.Printf("Error loading preferences: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded preferences: Theme=%s, Language=%s, Timezone=%s, Notifications=%v\n",
		prefs.Theme, prefs.Language, prefs.Timezone, prefs.NotificationsEnabled)
}