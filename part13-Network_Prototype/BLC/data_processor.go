package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) (bool, error) {
	var js interface{}
	err := json.Unmarshal(data, &js)
	if err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}
	return true, nil
}

// ParseUserData attempts to parse JSON bytes into a User struct.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ParseUserData(rawData []byte) (*User, error) {
	valid, err := ValidateJSON(rawData)
	if !valid {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(rawData, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}
	return &user, nil
}

func main() {
	jsonStr := `{"id": 1, "name": "Alice", "email": "alice@example.com"}`
	user, err := ParseUserData([]byte(jsonStr))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Parsed User: %+v\n", user)
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
		return errors.New("record value cannot be negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord, multiplier float64) DataRecord {
	return DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp.UTC(),
		Tags:      append(record.Tags, "processed"),
	}
}

func ProcessRecords(records []DataRecord, multiplier float64) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record, multiplier))
	}
	return processed, nil
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}
	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	return sum / float64(len(records))
}

func FilterByTag(records []DataRecord, tag string) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		for _, t := range record.Tags {
			if t == tag {
				filtered = append(filtered, record)
				break
			}
		}
	}
	return filtered
}
package main

import (
	"regexp"
	"strings"
	"unicode"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	cleaned := strings.TrimSpace(input)
	cleaned = dp.whitespaceRegex.ReplaceAllString(cleaned, " ")
	return cleaned
}

func (dp *DataProcessor) NormalizeWhitespace(input string) string {
	return dp.whitespaceRegex.ReplaceAllString(input, " ")
}

func (dp *DataProcessor) RemoveNonPrintable(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, input)
}

func (dp *DataProcessor) Process(input string) string {
	step1 := dp.RemoveNonPrintable(input)
	step2 := dp.NormalizeWhitespace(step1)
	return dp.CleanString(step2)
}