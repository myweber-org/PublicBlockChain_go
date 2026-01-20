package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

func ParseAndValidateJSON(rawData []byte, requiredFields []string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var missingFields []string
	for _, field := range requiredFields {
		if _, exists := data[field]; !exists {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		return nil, ValidationError{
			Field:   "required_fields",
			Message: fmt.Sprintf("missing required fields: %s", strings.Join(missingFields, ", ")),
		}
	}

	for key, value := range data {
		if strVal, ok := value.(string); ok && strings.TrimSpace(strVal) == "" {
			return nil, ValidationError{
				Field:   key,
				Message: "field cannot be empty",
			}
		}
	}

	return data, nil
}

func main() {
	jsonData := []byte(`{"name": "test", "age": 25, "email": ""}`)
	required := []string{"name", "age", "email"}

	result, err := ParseAndValidateJSON(jsonData, required)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Validated data: %v\n", result)
}