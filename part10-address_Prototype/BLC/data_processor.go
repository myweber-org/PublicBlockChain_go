
package main

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")
	
	// Trim spaces from beginning and end
	cleaned = strings.TrimSpace(cleaned)
	
	// Convert to lowercase for consistency
	cleaned = strings.ToLower(cleaned)
	
	return cleaned
}

func NormalizeString(input string) string {
	cleaned := CleanInput(input)
	
	// Remove special characters except alphanumeric and spaces
	re := regexp.MustCompile(`[^a-z0-9\s]`)
	normalized := re.ReplaceAllString(cleaned, "")
	
	return normalized
}

func ProcessData(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		processed := NormalizeString(input)
		if processed != "" {
			results = append(results, processed)
		}
	}
	return results
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
		return DataRecord{}, err
	}

	transformed := DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp.UTC(),
		Tags:      append([]string{}, record.Tags...),
	}

	if len(transformed.Tags) == 0 {
		transformed.Tags = []string{"default"}
	}

	return transformed, nil
}

func ProcessBatch(records []DataRecord, multiplier float64) ([]DataRecord, error) {
	var results []DataRecord
	var errors []string

	for i, record := range records {
		transformed, err := TransformRecord(record, multiplier)
		if err != nil {
			errors = append(errors, fmt.Sprintf("record %d: %v", i, err))
			continue
		}
		results = append(results, transformed)
	}

	if len(errors) > 0 {
		return results, fmt.Errorf("processing completed with errors: %s", strings.Join(errors, "; "))
	}

	return results, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, error) {
	if len(records) == 0 {
		return 0, 0, errors.New("no records provided for statistics calculation")
	}

	var sum float64
	for _, record := range records {
		sum += record.Value
	}

	mean := sum / float64(len(records))

	var varianceSum float64
	for _, record := range records {
		diff := record.Value - mean
		varianceSum += diff * diff
	}
	variance := varianceSum / float64(len(records))

	return mean, variance, nil
}
package main

import (
	"fmt"
	"strings"
	"unicode"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func NormalizeUsername(username string) string {
	trimmed := strings.TrimSpace(username)
	var result strings.Builder
	for _, r := range trimmed {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	localPart, domainPart, found := strings.Cut(email, "@")
	if !found || localPart == "" || domainPart == "" {
		return false
	}
	if strings.Contains(domainPart, "..") || strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") {
		return false
	}
	return true
}

func ProcessUserInput(username, email string, age int) (*UserData, error) {
	normalizedUsername := NormalizeUsername(username)
	if normalizedUsername == "" {
		return nil, fmt.Errorf("invalid username: contains no valid characters")
	}

	if !ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	if age < 0 || age > 150 {
		return nil, fmt.Errorf("age must be between 0 and 150")
	}

	return &UserData{
		Username: normalizedUsername,
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Age:      age,
	}, nil
}

func main() {
	user, err := ProcessUserInput("  John_Doe-123  ", "john@example.com", 30)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}