package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Comments string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	re := regexp.MustCompile(`[<>"'&]`)
	return re.ReplaceAllString(input, "")
}

func ValidateUserData(data UserData) (bool, []string) {
	var errors []string

	if len(data.Username) < 3 || len(data.Username) > 20 {
		errors = append(errors, "Username must be between 3 and 20 characters")
	}

	if !emailRegex.MatchString(data.Email) {
		errors = append(errors, "Invalid email format")
	}

	if len(data.Comments) > 500 {
		errors = append(errors, "Comments cannot exceed 500 characters")
	}

	return len(errors) == 0, errors
}

func ProcessUserData(data UserData) UserData {
	return UserData{
		Username: SanitizeInput(data.Username),
		Email:    strings.ToLower(SanitizeInput(data.Email)),
		Comments: SanitizeInput(data.Comments),
	}
}package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateAndTransform(data UserData) (UserData, error) {
	if strings.TrimSpace(data.Username) == "" {
		return UserData{}, fmt.Errorf("username cannot be empty")
	}

	if !strings.Contains(data.Email, "@") {
		return UserData{}, fmt.Errorf("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return UserData{}, fmt.Errorf("age must be between 0 and 150")
	}

	transformed := UserData{
		Username: strings.ToLower(strings.TrimSpace(data.Username)),
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Age:      data.Age,
	}

	return transformed, nil
}

func main() {
	sampleData := UserData{
		Username: "  TestUser  ",
		Email:    "TEST@EXAMPLE.COM",
		Age:      25,
	}

	result, err := ValidateAndTransform(sampleData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Original: %+v\n", sampleData)
	fmt.Printf("Processed: %+v\n", result)
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
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("invalid email format")
	}
	return email, nil
}

func validateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}
	pattern := `^[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(pattern, username)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	return nil
}

func processUserData(rawData []byte) (*UserData, error) {
	var data UserData
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	normalizedEmail, err := normalizeEmail(data.Email)
	if err != nil {
		return nil, fmt.Errorf("email validation failed: %v", err)
	}
	data.Email = normalizedEmail

	if err := validateUsername(data.Username); err != nil {
		return nil, fmt.Errorf("username validation failed: %v", err)
	}

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age must be between 0 and 150")
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email": "Test.User@Example.COM", "username": "valid_user123", "age": 25}`
	processedData, err := processUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processedData)
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

func validateUsername(username string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]{3,20}$", username)
	return matched
}

func validateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}

func transformData(data UserData) (UserData, error) {
	if !validateUsername(data.Username) {
		return UserData{}, fmt.Errorf("invalid username format")
	}

	if !validateEmail(data.Email) {
		return UserData{}, fmt.Errorf("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return UserData{}, fmt.Errorf("age must be between 0 and 150")
	}

	transformed := UserData{
		Username: strings.ToLower(data.Username),
		Email:    strings.ToLower(data.Email),
		Age:      data.Age,
	}

	return transformed, nil
}

func processJSONInput(jsonStr string) (string, error) {
	var userData UserData
	err := json.Unmarshal([]byte(jsonStr), &userData)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	transformed, err := transformData(userData)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(transformed, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return string(result), nil
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID      string
	Name    string
	Value   string
	IsValid bool
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Name:  strings.TrimSpace(row[1]),
			Value: strings.TrimSpace(row[2]),
		}

		record.IsValid = validateRecord(record)
		records = append(records, record)
	}

	return records, nil
}

func validateRecord(record DataRecord) bool {
	if record.ID == "" || record.Name == "" {
		return false
	}

	if len(record.Value) > 100 {
		return false
	}

	return true
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if record.IsValid {
			valid = append(valid, record)
		}
	}
	return valid
}

func GenerateSummary(records []DataRecord) {
	validCount := 0
	for _, record := range records {
		if record.IsValid {
			validCount++
		}
	}

	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Valid records: %d\n", validCount)
	fmt.Printf("Invalid records: %d\n", len(records)-validCount)
}
package data

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
	ErrEmptyField   = errors.New("required field is empty")
)

type DataRecord struct {
	ID        string
	Timestamp time.Time
	Value     float64
	Tags      []string
	Validated bool
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return ErrEmptyField
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if record.Value < 0 {
		return errors.New("value cannot be negative")
	}
	return nil
}

func NormalizeTags(tags []string) []string {
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

func TransformValue(value float64, multiplier float64) (float64, error) {
	if multiplier <= 0 {
		return 0, errors.New("multiplier must be positive")
	}
	return value * multiplier, nil
}

func ProcessRecord(record DataRecord, multiplier float64) (DataRecord, error) {
	if err := ValidateRecord(record); err != nil {
		return DataRecord{}, err
	}
	
	transformedValue, err := TransformValue(record.Value, multiplier)
	if err != nil {
		return DataRecord{}, err
	}
	
	normalizedTags := NormalizeTags(record.Tags)
	
	return DataRecord{
		ID:        record.ID,
		Timestamp: record.Timestamp,
		Value:     transformedValue,
		Tags:      normalizedTags,
		Validated: true,
	}, nil
}

func BatchProcess(records []DataRecord, multiplier float64) ([]DataRecord, []error) {
	var processed []DataRecord
	var errs []error
	
	for i, record := range records {
		processedRecord, err := ProcessRecord(record, multiplier)
		if err != nil {
			errs = append(errs, errors.New("record "+record.ID+" at index "+string(rune(i))+" failed: "+err.Error()))
			continue
		}
		processed = append(processed, processedRecord)
	}
	
	return processed, errs
}