package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type UserPreferences struct {
	Theme       string `json:"theme"`
	Language    string `json:"language"`
	Notifications bool `json:"notifications"`
	Timezone    string `json:"timezone"`
}

func LoadPreferences(filename string) (*UserPreferences, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open preferences file: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read preferences file: %w", err)
	}

	var prefs UserPreferences
	err = json.Unmarshal(bytes, &prefs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse preferences JSON: %w", err)
	}

	return &prefs, nil
}

func SavePreferences(filename string, prefs *UserPreferences) error {
	bytes, err := json.MarshalIndent(prefs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %w", err)
	}

	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write preferences file: %w", err)
	}

	return nil
}

func main() {
	prefs, err := LoadPreferences("user_prefs.json")
	if err != nil {
		fmt.Printf("Error loading preferences: %v\n", err)
		return
	}

	fmt.Printf("Loaded preferences: %+v\n", prefs)

	newPrefs := &UserPreferences{
		Theme:       "dark",
		Language:    "en-US",
		Notifications: true,
		Timezone:    "UTC",
	}

	err = SavePreferences("user_prefs_updated.json", newPrefs)
	if err != nil {
		fmt.Printf("Error saving preferences: %v\n", err)
		return
	}

	fmt.Println("Preferences saved successfully")
}