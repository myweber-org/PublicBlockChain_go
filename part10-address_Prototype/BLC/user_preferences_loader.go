package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type UserPreferences struct {
	Theme      string `json:"theme"`
	Language   string `json:"language"`
	Timezone   string `json:"timezone"`
	DateFormat string `json:"date_format"`
}

type PreferencesCache struct {
	mu          sync.RWMutex
	preferences map[string]UserPreferences
	expiry      map[string]time.Time
	ttl         time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
	return &PreferencesCache{
		preferences: make(map[string]UserPreferences),
		expiry:      make(map[string]time.Time),
		ttl:         ttl,
	}
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	prefs, exists := c.preferences[userID]
	if !exists {
		return UserPreferences{}, false
	}

	expiry, expiryExists := c.expiry[userID]
	if !expiryExists || time.Now().After(expiry) {
		return UserPreferences{}, false
	}

	return prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.preferences[userID] = prefs
	c.expiry[userID] = time.Now().Add(c.ttl)
}

func loadPreferencesFromFile(filename string) (map[string]UserPreferences, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read preferences file: %w", err)
	}

	var preferences map[string]UserPreferences
	if err := json.Unmarshal(data, &preferences); err != nil {
		return nil, fmt.Errorf("failed to parse preferences JSON: %w", err)
	}

	return preferences, nil
}

func validatePreferences(prefs UserPreferences) error {
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
		return fmt.Errorf("date_format cannot be empty")
	}

	validThemes := map[string]bool{"light": true, "dark": true, "auto": true}
	if !validThemes[prefs.Theme] {
		return fmt.Errorf("invalid theme: %s", prefs.Theme)
	}

	return nil
}

func LoadUserPreferences(userID string, cache *PreferencesCache, filename string) (UserPreferences, error) {
	if cache != nil {
		if prefs, found := cache.Get(userID); found {
			return prefs, nil
		}
	}

	allPrefs, err := loadPreferencesFromFile(filename)
	if err != nil {
		return UserPreferences{}, err
	}

	prefs, exists := allPrefs[userID]
	if !exists {
		return UserPreferences{}, fmt.Errorf("preferences not found for user: %s", userID)
	}

	if err := validatePreferences(prefs); err != nil {
		return UserPreferences{}, fmt.Errorf("invalid preferences for user %s: %w", userID, err)
	}

	if cache != nil {
		cache.Set(userID, prefs)
	}

	return prefs, nil
}

func main() {
	cache := NewPreferencesCache(5 * time.Minute)

	prefs, err := LoadUserPreferences("user123", cache, "preferences.json")
	if err != nil {
		fmt.Printf("Error loading preferences: %v\n", err)
		return
	}

	fmt.Printf("Loaded preferences: %+v\n", prefs)
}