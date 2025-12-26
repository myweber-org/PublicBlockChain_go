package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

type UserPreferences struct {
    Theme        string `json:"theme"`
    Notifications bool   `json:"notifications"`
    Language     string `json:"language"`
    ItemsPerPage int    `json:"items_per_page"`
}

func LoadPreferences(filename string) (*UserPreferences, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open preferences file: %w", err)
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read preferences file: %w", err)
    }

    var prefs UserPreferences
    if err := json.Unmarshal(data, &prefs); err != nil {
        return nil, fmt.Errorf("failed to parse preferences JSON: %w", err)
    }

    if err := validatePreferences(&prefs); err != nil {
        return nil, fmt.Errorf("preferences validation failed: %w", err)
    }

    return &prefs, nil
}

func validatePreferences(prefs *UserPreferences) error {
    if prefs.Theme != "light" && prefs.Theme != "dark" {
        return fmt.Errorf("invalid theme value: %s", prefs.Theme)
    }
    if prefs.ItemsPerPage < 5 || prefs.ItemsPerPage > 100 {
        return fmt.Errorf("items per page must be between 5 and 100, got %d", prefs.ItemsPerPage)
    }
    if prefs.Language == "" {
        return fmt.Errorf("language cannot be empty")
    }
    return nil
}

func main() {
    prefs, err := LoadPreferences("user_prefs.json")
    if err != nil {
        fmt.Printf("Error loading preferences: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Loaded preferences: %+v\n", prefs)
}