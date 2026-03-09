
package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", ErrEmptyInput
	}

	pattern := `^[a-zA-Z0-9_\-\.]+$`
	matched, err := regexp.MatchString(pattern, trimmed)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", ErrInvalidCharacters
	}

	if len(trimmed) > 50 {
		return "", ErrInputTooLong
	}
	return trimmed, nil
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

var (
	ErrEmptyInput        = errors.New("input cannot be empty")
	ErrInvalidCharacters = errors.New("input contains invalid characters")
	ErrInputTooLong      = errors.New("input exceeds maximum length")
)
package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

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

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
        }

        name := row[1]
        if name == "" {
            return nil, fmt.Errorf("empty name at line %d", lineNumber)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
        }

        records = append(records, DataRecord{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    min := records[0].Value
    max := records[0].Value

    for _, record := range records {
        sum += record.Value
        if record.Value < min {
            min = record.Value
        }
        if record.Value > max {
            max = record.Value
        }
    }

    average := sum / float64(len(records))
    return average, max - min, len(records)
}

func ValidateRecord(record DataRecord) error {
    if record.ID <= 0 {
        return errors.New("ID must be positive integer")
    }
    if record.Name == "" {
        return errors.New("name cannot be empty")
    }
    if record.Value < 0 {
        return errors.New("value cannot be negative")
    }
    return nil
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

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func sanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func processUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !validateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = sanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email":"test@example.com","username":"  john_doe  ","age":25}`
	processedData, err := processUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processedData)
}package main

import (
	"errors"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func NormalizeUserData(data UserData) UserData {
	return UserData{
		Username: strings.ToLower(strings.TrimSpace(data.Username)),
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Age:      data.Age,
	}
}

func ProcessUserInput(rawData UserData) (UserData, error) {
	if err := ValidateUserData(rawData); err != nil {
		return UserData{}, err
	}
	return NormalizeUserData(rawData), nil
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
	ID    string
	Name  string
	Email string
	Valid bool
}

func processCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		if len(line) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(line[0]),
			Name:  strings.TrimSpace(line[1]),
			Email: strings.TrimSpace(line[2]),
			Valid: validateRecord(strings.TrimSpace(line[0]), strings.TrimSpace(line[2])),
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(id, email string) bool {
	if id == "" || email == "" {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func generateReport(records []DataRecord) {
	validCount := 0
	invalidCount := 0

	fmt.Println("Data Processing Report")
	fmt.Println("======================")

	for _, record := range records {
		if record.Valid {
			validCount++
			fmt.Printf("✓ Valid: %s - %s\n", record.ID, record.Name)
		} else {
			invalidCount++
			fmt.Printf("✗ Invalid: %s - %s (Email: %s)\n", record.ID, record.Name, record.Email)
		}
	}

	fmt.Printf("\nSummary: %d valid, %d invalid records\n", validCount, invalidCount)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run data_processor.go <csv_file>")
		return
	}

	filename := os.Args[1]
	records, err := processCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	generateReport(records)
}
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (dp *DataProcessor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return normalized
}

func (dp *DataProcessor) NormalizeCase(input string, toUpper bool) string {
	cleaned := dp.CleanInput(input)
	if toUpper {
		return strings.ToUpper(cleaned)
	}
	return strings.ToLower(cleaned)
}

func (dp *DataProcessor) ExtractAlphanumeric(input string) string {
	alnumRegex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return alnumRegex.ReplaceAllString(input, "")
}

func main() {
	processor := NewDataProcessor()
	
	sample := "  Hello   World! 123  "
	cleaned := processor.CleanInput(sample)
	normalized := processor.NormalizeCase(cleaned, false)
	alnum := processor.ExtractAlphanumeric(sample)
	
	println("Original:", sample)
	println("Cleaned:", cleaned)
	println("Normalized:", normalized)
	println("Alphanumeric only:", alnum)
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

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Value = record.Value * 1.1
	transformed.Tags = append(record.Tags, "processed")
	transformed.Tags = normalizeTags(transformed.Tags)
	return transformed
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

func ProcessData(records []DataRecord) ([]DataRecord, error) {
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
		{
			ID:        "rec001",
			Value:     100.0,
			Timestamp: time.Now(),
			Tags:      []string{"test", "sample"},
		},
		{
			ID:        "rec002",
			Value:     200.0,
			Timestamp: time.Now().Add(-time.Hour),
			Tags:      []string{"production", "SAMPLE"},
		},
	}

	processed, err := ProcessData(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	for _, rec := range processed {
		fmt.Printf("Processed: %s - %.2f - %v\n", rec.ID, rec.Value, rec.Tags)
	}
}
package main

import (
    "regexp"
    "strings"
)

type DataCleaner struct {
    spacesRegex *regexp.Regexp
    emailRegex  *regexp.Regexp
}

func NewDataCleaner() *DataCleaner {
    return &DataCleaner{
        spacesRegex: regexp.MustCompile(`\s+`),
        emailRegex:  regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
    }
}

func (dc *DataCleaner) TrimSpaces(input string) string {
    return strings.TrimSpace(input)
}

func (dc *DataCleaner) NormalizeSpaces(input string) string {
    trimmed := dc.TrimSpaces(input)
    return dc.spacesRegex.ReplaceAllString(trimmed, " ")
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
    return dc.emailRegex.MatchString(email)
}

func (dc *DataCleaner) ProcessUserInput(rawInput string) (string, bool) {
    normalized := dc.NormalizeSpaces(rawInput)
    
    if normalized == "" {
        return "", false
    }
    
    return normalized, true
}

func main() {
    cleaner := NewDataCleaner()
    
    testInputs := []string{
        "  hello   world  ",
        "user@example.com",
        "  multiple   spaces   here  ",
        "",
    }
    
    for _, input := range testInputs {
        processed, valid := cleaner.ProcessUserInput(input)
        if valid {
            println("Processed:", processed)
        } else {
            println("Invalid input:", input)
        }
    }
}