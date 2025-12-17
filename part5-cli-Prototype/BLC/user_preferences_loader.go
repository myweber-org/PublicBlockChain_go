package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type UserPreferences struct {
    Theme     string `json:"theme"`
    Language  string `json:"language"`
    Notifications bool `json:"notifications"`
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
        return nil, fmt.Errorf("failed to decode preferences: %w", err)
    }

    return &prefs, nil
}

func main() {
    prefs, err := LoadPreferences("preferences.json")
    if err != nil {
        fmt.Printf("Error loading preferences: %v\n", err)
        return
    }

    fmt.Printf("Loaded preferences: Theme=%s, Language=%s, Notifications=%v\n",
        prefs.Theme, prefs.Language, prefs.Notifications)
}