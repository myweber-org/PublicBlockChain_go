package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) (bool, error) {
	var js interface{}
	err := json.Unmarshal(data, &js)
	if err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}
	return true, nil
}

// ParseUserData attempts to parse JSON data into a map.
func ParseUserData(jsonData []byte) (map[string]interface{}, error) {
	valid, err := ValidateJSON(jsonData)
	if !valid {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func main() {
	sampleJSON := []byte(`{"name": "Alice", "age": 30, "active": true}`)

	parsed, err := ParseUserData(sampleJSON)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Parsed data: %v\n", parsed)
}