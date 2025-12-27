
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) bool {
	var js interface{}
	return json.Unmarshal(data, &js) == nil
}

// ParseJSONMap attempts to parse the byte slice into a map[string]interface{}.
func ParseJSONMap(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func main() {
	// Example usage
	validJSON := []byte(`{"name": "test", "value": 42}`)
	invalidJSON := []byte(`{name: test}`)

	fmt.Println("Valid JSON check:", ValidateJSON(validJSON))
	fmt.Println("Invalid JSON check:", ValidateJSON(invalidJSON))

	parsed, err := ParseJSONMap(validJSON)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed data: %v\n", parsed)
}