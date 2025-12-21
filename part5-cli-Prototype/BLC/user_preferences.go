package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserPreferences struct {
	Theme        string  `json:"theme"`
	Notifications bool   `json:"notifications"`
	Language     string  `json:"language"`
	Volume       float64 `json:"volume"`
}

func (up *UserPreferences) Validate() error {
	if up.Theme == "" {
		up.Theme = "light"
	}
	if up.Theme != "light" && up.Theme != "dark" {
		return fmt.Errorf("invalid theme: %s", up.Theme)
	}
	if up.Language == "" {
		up.Language = "en"
	}
	if up.Volume < 0.0 || up.Volume > 1.0 {
		return fmt.Errorf("volume must be between 0.0 and 1.0, got: %f", up.Volume)
	}
	return nil
}

func main() {
	prefs := UserPreferences{
		Theme:    "dark",
		Language: "es",
		Volume:   0.8,
	}
	
	if err := prefs.Validate(); err != nil {
		log.Fatal(err)
	}
	
	data, err := json.MarshalIndent(prefs, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Validated preferences:")
	fmt.Println(string(data))
}