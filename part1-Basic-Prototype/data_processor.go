package main

import (
	"errors"
	"strings"
)

type DataRecord struct {
	ID    string
	Value string
	Valid bool
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if len(record.Value) > 100 {
		return errors.New("value exceeds maximum length")
	}
	return nil
}

func TransformValue(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, rec := range records {
		if err := ValidateRecord(rec); err != nil {
			return nil, err
		}
		rec.Value = TransformValue(rec.Value)
		rec.Valid = true
		processed = append(processed, rec)
	}
	return processed, nil
}
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateAndParseJSON(input []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(input, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if data.ID <= 0 {
		return nil, fmt.Errorf("invalid ID: must be positive integer")
	}
	if data.Name == "" {
		return nil, fmt.Errorf("name field cannot be empty")
	}
	if data.Email == "" {
		return nil, fmt.Errorf("email field cannot be empty")
	}

	return &data, nil
}

func main() {
	jsonInput := `{"id": 123, "name": "John Doe", "email": "john@example.com"}`
	parsedData, err := ValidateAndParseJSON([]byte(jsonInput))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Parsed data: %+v\n", parsedData)
}