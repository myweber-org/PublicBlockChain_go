
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
}

func ValidateUsername(username string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]{3,20}$", username)
	return matched
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

func TransformProfile(profile UserProfile) (UserProfile, error) {
	if !ValidateUsername(profile.Username) {
		return profile, fmt.Errorf("invalid username format")
	}

	if !ValidateEmail(profile.Email) {
		return profile, fmt.Errorf("invalid email format")
	}

	if profile.Age < 0 || profile.Age > 150 {
		return profile, fmt.Errorf("age must be between 0 and 150")
	}

	profile.Username = strings.TrimSpace(profile.Username)
	profile.Email = strings.ToLower(strings.TrimSpace(profile.Email))

	return profile, nil
}

func ProcessUserData(jsonData []byte) (UserProfile, error) {
	var profile UserProfile
	err := json.Unmarshal(jsonData, &profile)
	if err != nil {
		return profile, fmt.Errorf("failed to parse JSON: %v", err)
	}

	transformedProfile, err := TransformProfile(profile)
	if err != nil {
		return profile, fmt.Errorf("validation failed: %v", err)
	}

	return transformedProfile, nil
}

func main() {
	jsonInput := `{"id":1,"username":"john_doe","email":"John@Example.COM","age":25,"active":true}`

	processedProfile, err := ProcessUserData([]byte(jsonInput))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	output, _ := json.MarshalIndent(processedProfile, "", "  ")
	fmt.Printf("Processed Profile:\n%s\n", output)
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
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord, multiplier float64) (DataRecord, error) {
	if err := ValidateRecord(record); err != nil {
		return DataRecord{}, fmt.Errorf("validation failed: %w", err)
	}

	transformed := DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp,
		Tags:      append([]string{}, record.Tags...),
	}

	transformed.Tags = append(transformed.Tags, "processed")

	return transformed, nil
}

func ProcessBatch(records []DataRecord, multiplier float64) ([]DataRecord, []error) {
	var processed []DataRecord
	var errs []error

	for i, record := range records {
		transformed, err := TransformRecord(record, multiplier)
		if err != nil {
			errs = append(errs, fmt.Errorf("record %d: %w", i, err))
			continue
		}
		processed = append(processed, transformed)
	}

	return processed, errs
}

func CalculateStatistics(records []DataRecord) (float64, float64, error) {
	if len(records) == 0 {
		return 0, 0, errors.New("no records provided")
	}

	var sum float64
	var count int

	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			continue
		}
		sum += record.Value
		count++
	}

	if count == 0 {
		return 0, 0, errors.New("no valid records found")
	}

	average := sum / float64(count)
	return sum, average, nil
}