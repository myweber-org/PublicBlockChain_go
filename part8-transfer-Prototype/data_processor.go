
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
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord) (DataRecord, error) {
	if err := ValidateRecord(record); err != nil {
		return DataRecord{}, err
	}

	transformed := record
	transformed.Value = record.Value * 1.1
	transformed.Tags = append(record.Tags, "processed")
	transformed.Tags = normalizeTags(transformed.Tags)
	return transformed, nil
}

func normalizeTags(tags []string) []string {
	uniqueTags := make(map[string]bool)
	var result []string

	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" && !uniqueTags[normalized] {
			uniqueTags[normalized] = true
			result = append(result, normalized)
		}
	}
	return result
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	var errors []string

	for i, record := range records {
		transformed, err := TransformRecord(record)
		if err != nil {
			errors = append(errors, fmt.Sprintf("record %d: %v", i, err))
			continue
		}
		processed = append(processed, transformed)
	}

	if len(errors) > 0 {
		return processed, fmt.Errorf("processing errors: %s", strings.Join(errors, "; "))
	}
	return processed, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, error) {
	if len(records) == 0 {
		return 0, 0, errors.New("no records to process")
	}

	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	average := sum / float64(len(records))

	var varianceSum float64
	for _, record := range records {
		diff := record.Value - average
		varianceSum += diff * diff
	}
	variance := varianceSum / float64(len(records))

	return average, variance, nil
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func normalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return "", fmt.Errorf("regex error: %v", err)
	}
	if !matched {
		return "", fmt.Errorf("invalid email format")
	}
	return email, nil
}

func validateUsername(username string) (string, error) {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return "", fmt.Errorf("username must be between 3 and 20 characters")
	}
	pattern := `^[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(pattern, username)
	if err != nil {
		return "", fmt.Errorf("regex error: %v", err)
	}
	if !matched {
		return "", fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	return username, nil
}

func validateAge(age int) error {
	if age < 0 || age > 150 {
		return fmt.Errorf("age must be between 0 and 150")
	}
	return nil
}

func ProcessUserData(rawData []byte) (*UserData, error) {
	var data UserData
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	email, err := normalizeEmail(data.Email)
	if err != nil {
		return nil, fmt.Errorf("email validation failed: %v", err)
	}
	data.Email = email

	username, err := validateUsername(data.Username)
	if err != nil {
		return nil, fmt.Errorf("username validation failed: %v", err)
	}
	data.Username = username

	if err := validateAge(data.Age); err != nil {
		return nil, fmt.Errorf("age validation failed: %v", err)
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email": "  TEST@Example.COM ", "username": "user_123", "age": 25}`
	processedData, err := ProcessUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed Data: %+v\n", processedData)
}
package main

import (
	"fmt"
)

// FilterAndDouble filters even numbers from a slice and doubles their values.
func FilterAndDouble(numbers []int) []int {
	var result []int
	for _, num := range numbers {
		if num%2 == 0 {
			result = append(result, num*2)
		}
	}
	return result
}

func main() {
	input := []int{1, 2, 3, 4, 5, 6}
	output := FilterAndDouble(input)
	fmt.Println("Original:", input)
	fmt.Println("Filtered and doubled:", output)
}