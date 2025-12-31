package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}
	if data.Age < 18 || data.Age > 120 {
		return errors.New("age must be between 18 and 120")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
	transformedUsername := TransformUsername(username)
	userData := UserData{
		Username: transformedUsername,
		Email:    strings.TrimSpace(email),
		Age:      age,
	}
	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}
	return userData, nil
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// DataPayload represents a simple JSON structure
type DataPayload struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
	Version string `json:"version,omitempty"`
}

// ValidatePayload checks if the DataPayload has valid fields
func ValidatePayload(p *DataPayload) error {
	if p.ID <= 0 {
		return fmt.Errorf("invalid ID: must be positive integer")
	}
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}

// ParseJSONData unmarshals JSON bytes into a DataPayload and validates it
func ParseJSONData(rawData []byte) (*DataPayload, error) {
	var payload DataPayload
	if err := json.Unmarshal(rawData, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if err := ValidatePayload(&payload); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &payload, nil
}

func main() {
	// Example JSON data
	jsonStr := `{"id": 101, "name": "TestItem", "active": true}`

	payload, err := ParseJSONData([]byte(jsonStr))
	if err != nil {
		log.Fatalf("Error processing data: %v", err)
	}

	fmt.Printf("Processed payload: ID=%d, Name=%s, Active=%t\n",
		payload.ID, payload.Name, payload.Active)
}