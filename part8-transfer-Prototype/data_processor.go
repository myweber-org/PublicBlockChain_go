
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

// ParseJSONToMap parses JSON bytes into a map[string]interface{}.
func ParseJSONToMap(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func main() {
	sampleJSON := []byte(`{"name": "test", "value": 42}`)

	valid, err := ValidateJSON(sampleJSON)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}
	fmt.Printf("JSON is valid: %v\n", valid)

	parsed, err := ParseJSONToMap(sampleJSON)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}
	fmt.Printf("Parsed data: %v\n", parsed)
}