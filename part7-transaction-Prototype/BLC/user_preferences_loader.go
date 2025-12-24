package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "sync"
    "time"
)

type UserPreferences struct {
    Theme      string   `json:"theme"`
    Language   string   `json:"language"`
    Notifications bool  `json:"notifications"`
    Timezone   string   `json:"timezone"`
}

type PreferencesCache struct {
    mu      sync.RWMutex
    data    map[string]UserPreferences
    ttl     time.Duration
    updated map[string]time.Time
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
    return &PreferencesCache{
        data:    make(map[string]UserPreferences),
        updated: make(map[string]time.Time),
        ttl:     ttl,
    }
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    prefs, exists := c.data[userID]
    if !exists {
        return UserPreferences{}, false
    }

    if time.Since(c.updated[userID]) > c.ttl {
        return UserPreferences{}, false
    }

    return prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.data[userID] = prefs
    c.updated[userID] = time.Now()
}

func validatePreferences(prefs UserPreferences) error {
    validThemes := map[string]bool{"light": true, "dark": true, "auto": true}
    if !validThemes[prefs.Theme] {
        return errors.New("invalid theme selection")
    }

    if prefs.Language == "" {
        return errors.New("language cannot be empty")
    }

    if prefs.Timezone == "" {
        return errors.New("timezone cannot be empty")
    }

    return nil
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

    for userID, prefs := range preferences {
        if err := validatePreferences(prefs); err != nil {
            return nil, fmt.Errorf("invalid preferences for user %s: %w", userID, err)
        }
    }

    return preferences, nil
}

func main() {
    cache := NewPreferencesCache(5 * time.Minute)

    preferences, err := loadPreferencesFromFile("preferences.json")
    if err != nil {
        fmt.Printf("Error loading preferences: %v\n", err)
        return
    }

    for userID, prefs := range preferences {
        cache.Set(userID, prefs)
        fmt.Printf("Loaded preferences for user: %s\n", userID)
    }

    if prefs, found := cache.Get("user123"); found {
        fmt.Printf("Cached preferences: %+v\n", prefs)
    }
}