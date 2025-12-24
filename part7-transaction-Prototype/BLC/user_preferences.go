package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserPreferences struct {
	Theme       string  `json:"theme"`
	Language    string  `json:"language"`
	VolumeLevel float64 `json:"volume_level"`
	NotificationsEnabled bool `json:"notifications_enabled"`
}

func (up *UserPreferences) Validate() error {
	if up.Theme == "" {
		up.Theme = "light"
	}
	if up.Language == "" {
		up.Language = "en"
	}
	if up.VolumeLevel < 0.0 || up.VolumeLevel > 1.0 {
		return fmt.Errorf("volume level must be between 0.0 and 1.0, got %f", up.VolumeLevel)
	}
	return nil
}

func (up *UserPreferences) SetDefaults() {
	if up.Theme == "" {
		up.Theme = "light"
	}
	if up.Language == "" {
		up.Language = "en"
	}
	if up.VolumeLevel == 0.0 {
		up.VolumeLevel = 0.5
	}
}

func main() {
	prefJSON := `{"theme":"dark","volume_level":0.8,"notifications_enabled":true}`
	var prefs UserPreferences
	err := json.Unmarshal([]byte(prefJSON), &prefs)
	if err != nil {
		log.Fatal(err)
	}

	prefs.SetDefaults()
	err = prefs.Validate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Validated preferences: %+v\n", prefs)
}