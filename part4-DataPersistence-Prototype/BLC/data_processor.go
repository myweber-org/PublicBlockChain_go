
package data_processor

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

type DataRecord struct {
	ID        string  `json:"id"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
	Category  string  `json:"category"`
}

func ParseAndValidateJSON(rawData []byte) (*DataRecord, error) {
	var record DataRecord
	
	if err := json.Unmarshal(rawData, &record); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var validationErrors []ValidationError

	if strings.TrimSpace(record.ID) == "" {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "id",
			Message: "cannot be empty",
		})
	}

	if record.Value < 0 {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "value",
			Message: "must be non-negative",
		})
	}

	if record.Timestamp <= 0 {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "timestamp",
			Message: "must be positive integer",
		})
	}

	if !isValidCategory(record.Category) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "category",
			Message: "invalid category specified",
		})
	}

	if len(validationErrors) > 0 {
		var errorMessages []string
		for _, err := range validationErrors {
			errorMessages = append(errorMessages, err.Error())
		}
		return nil, fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
	}

	return &record, nil
}

func isValidCategory(category string) bool {
	validCategories := map[string]bool{
		"standard": true,
		"premium":  true,
		"legacy":   true,
	}
	return validCategories[category]
}
package main

import (
	"strings"
	"unicode"
)

func ProcessInput(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	trimmed := strings.TrimSpace(input)
	var cleaned strings.Builder

	for _, r := range trimmed {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			cleaned.WriteRune(r)
		}
	}

	result := strings.Join(strings.Fields(cleaned.String()), " ")
	return result, nil
}