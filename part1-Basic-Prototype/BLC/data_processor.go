
package main

import (
	"errors"
	"strings"
	"unicode"
)

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return errors.New("username can only contain letters, digits, underscores, and hyphens")
		}
	}

	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func SanitizeInput(input string) string {
	replacer := strings.NewReplacer(
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(input)
}

func ValidatePasswordStrength(password string) (bool, []string) {
	var issues []string

	if len(password) < 8 {
		issues = append(issues, "password must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if !hasUpper {
		issues = append(issues, "password must contain at least one uppercase letter")
	}
	if !hasLower {
		issues = append(issues, "password must contain at least one lowercase letter")
	}
	if !hasDigit {
		issues = append(issues, "password must contain at least one digit")
	}
	if !hasSpecial {
		issues = append(issues, "password must contain at least one special character")
	}

	return len(issues) == 0, issues
}
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        string
	Email     string
	Username  string
	Age       int
	Active    bool
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateProfile(p UserProfile) error {
	if p.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if !emailRegex.MatchString(p.Email) {
		return errors.New("invalid email format")
	}
	if len(p.Username) < 3 || len(p.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if p.Age < 0 || p.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformProfile(p UserProfile) UserProfile {
	p.Username = strings.ToLower(strings.TrimSpace(p.Username))
	p.Email = strings.ToLower(strings.TrimSpace(p.Email))
	return p
}

func ProcessUserProfile(p UserProfile) (UserProfile, error) {
	if err := ValidateProfile(p); err != nil {
		return UserProfile{}, err
	}
	transformed := TransformProfile(p)
	return transformed, nil
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataProcessor struct {
	InputPath  string
	OutputPath string
}

func NewDataProcessor(input, output string) *DataProcessor {
	return &DataProcessor{
		InputPath:  input,
		OutputPath: output,
	}
}

func (dp *DataProcessor) Process() error {
	inputFile, err := os.Open(dp.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dp.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	cleanedHeaders := dp.cleanHeaders(headers)
	if err := writer.Write(cleanedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		cleanedRecord := dp.cleanRecord(record)
		if dp.isValidRecord(cleanedRecord) {
			if err := writer.Write(cleanedRecord); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
			recordCount++
		}
	}

	fmt.Printf("Processed %d valid records\n", recordCount)
	return nil
}

func (dp *DataProcessor) cleanHeaders(headers []string) []string {
	cleaned := make([]string, len(headers))
	for i, header := range headers {
		cleaned[i] = strings.TrimSpace(header)
		cleaned[i] = strings.ToLower(cleaned[i])
		cleaned[i] = strings.ReplaceAll(cleaned[i], " ", "_")
	}
	return cleaned
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
	cleaned := make([]string, len(record))
	for i, field := range record {
		cleaned[i] = strings.TrimSpace(field)
		if cleaned[i] == "" {
			cleaned[i] = "N/A"
		}
	}
	return cleaned
}

func (dp *DataProcessor) isValidRecord(record []string) bool {
	for _, field := range record {
		if field == "" {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	processor := NewDataProcessor(os.Args[1], os.Args[2])
	if err := processor.Process(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func parseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
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

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
		if err != nil {
			return nil, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
		}

		records = append(records, Record{
			ID:     id,
			Name:   name,
			Value:  value,
			Active: active,
		})
	}

	return records, nil
}

func validateRecords(records []Record) error {
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", record.ID)
		}
		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID %d found", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("negative value %f for record ID %d", record.Value, record.ID)
		}
	}
	return nil
}

func processData(filename string) error {
	records, err := parseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	if err := validateRecords(records); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	totalValue := 0.0
	activeCount := 0
	for _, record := range records {
		totalValue += record.Value
		if record.Active {
			activeCount++
		}
	}

	fmt.Printf("Processed %d records successfully\n", len(records))
	fmt.Printf("Total value: %.2f\n", totalValue)
	fmt.Printf("Active records: %d\n", activeCount)
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	if err := processData(os.Args[1]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}