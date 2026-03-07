package main

import (
	"errors"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

func ValidateUserData(data UserData) error {
	if data.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	if data.Username == "" {
		return errors.New("username is required")
	}
	if len(data.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	transformedUsername := TransformUsername(username)
	userData := UserData{
		Email:    strings.TrimSpace(email),
		Username: transformedUsername,
		Age:      age,
	}
	err := ValidateUserData(userData)
	return userData, err
}
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

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return errors.New("username can only contain letters, numbers and underscores")
	}
	return nil
}

func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ProcessUserData(data UserData) (UserData, error) {
	if err := ValidateUsername(data.Username); err != nil {
		return UserData{}, err
	}
	
	if err := ValidateEmail(data.Email); err != nil {
		return UserData{}, err
	}
	
	data.Email = NormalizeEmail(data.Email)
	
	if data.Age < 0 || data.Age > 120 {
		return UserData{}, errors.New("age must be between 0 and 120")
	}
	
	return data, nil
}

func TransformUsername(username string) string {
	return strings.ReplaceAll(username, "_", "-")
}
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Valid     bool
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.After(time.Now()) {
		return errors.New("timestamp cannot be in the future")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.ID = strings.ToUpper(record.ID)
	transformed.Value = record.Value * 1.1
	return transformed
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record))
	}
	return processed, nil
}

func main() {
	records := []DataRecord{
		{ID: "abc123", Value: 100.0, Timestamp: time.Now().Add(-time.Hour), Valid: true},
		{ID: "def456", Value: 200.0, Timestamp: time.Now().Add(-2 * time.Hour), Valid: true},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	for _, record := range processed {
		fmt.Printf("Processed: ID=%s, Value=%.2f\n", record.ID, record.Value)
	}
}package main

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
	if len(data.Username) < 3 || len(data.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
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
	err := ValidateUserData(userData)
	if err != nil {
		return UserData{}, err
	}
	return userData, nil
}