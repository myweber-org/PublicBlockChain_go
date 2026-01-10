package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type UserPreferences struct {
	Theme      string `json:"theme"`
	Language   string `json:"language"`
	Timezone   string `json:"timezone"`
	DateFormat string `json:"date_format"`
}

type PreferencesCache struct {
	preferences *UserPreferences
	loadedAt    time.Time
	ttl         time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
	return &PreferencesCache{
		ttl: ttl,
	}
}

func (c *PreferencesCache) LoadFromFile(filename string) (*UserPreferences, error) {
	if c.preferences != nil && time.Since(c.loadedAt) < c.ttl {
		return c.preferences, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open preferences file: %w", err)
	}
	defer file.Close()

	var prefs UserPreferences
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&prefs); err != nil {
		return nil, fmt.Errorf("failed to decode preferences: %w", err)
	}

	if err := validatePreferences(&prefs); err != nil {
		return nil, fmt.Errorf("invalid preferences: %w", err)
	}

	c.preferences = &prefs
	c.loadedAt = time.Now()
	return c.preferences, nil
}

func validatePreferences(prefs *UserPreferences) error {
	if prefs.Theme == "" {
		return fmt.Errorf("theme cannot be empty")
	}
	if prefs.Language == "" {
		return fmt.Errorf("language cannot be empty")
	}
	if prefs.Timezone == "" {
		return fmt.Errorf("timezone cannot be empty")
	}
	if prefs.DateFormat == "" {
		return fmt.Errorf("date format cannot be empty")
	}
	return nil
}

func main() {
	cache := NewPreferencesCache(5 * time.Minute)
	prefs, err := cache.LoadFromFile("preferences.json")
	if err != nil {
		fmt.Printf("Error loading preferences: %v\n", err)
		return
	}
	fmt.Printf("Loaded preferences: %+v\n", prefs)
}