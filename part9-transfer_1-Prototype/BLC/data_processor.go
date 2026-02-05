package data

import (
	"errors"
	"regexp"
	"strings"
)

// ValidateEmail checks if the provided string is a valid email address.
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}
	return matched
}

// SanitizeInput removes leading and trailing whitespace from a string.
func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

// TransformToUpper converts all letters in the string to uppercase.
func TransformToUpper(input string) string {
	return strings.ToUpper(input)
}

// ParseInteger safely parses a string to an integer.
func ParseInteger(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0, errors.New("failed to parse integer")
	}
	return result, nil
}package main

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

// ParseUserData attempts to parse JSON data into a generic map.
func ParseUserData(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}
	return result, nil
}

func main() {
	sampleJSON := `{"name": "alice", "age": 30, "active": true}`

	valid, err := ValidateJSON([]byte(sampleJSON))
	if err != nil {
		log.Printf("Validation error: %v", err)
	} else {
		fmt.Println("JSON is valid:", valid)
	}

	userData, err := ParseUserData(sampleJSON)
	if err != nil {
		log.Printf("Parse error: %v", err)
	} else {
		fmt.Printf("Parsed data: %+v\n", userData)
	}
}
package main

import (
	"encoding/csv"
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

func ReadCSVFile(filename string) ([]DataRecord, error) {
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

		if lineNumber == 0 {
			lineNumber++
			continue
		}

		if len(line) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(line))
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := line[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
		lineNumber++
	}

	return records, nil
}

func CalculateStats(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var max float64
	count := len(records)

	for i, record := range records {
		sum += record.Value
		if i == 0 || record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(count)
	return average, max, count
}

func ValidateRecords(records []DataRecord) []string {
	errors := []string{}
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			errors = append(errors, fmt.Sprintf("invalid ID %d: must be positive", record.ID))
		}

		if seenIDs[record.ID] {
			errors = append(errors, fmt.Sprintf("duplicate ID %d found", record.ID))
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			errors = append(errors, fmt.Sprintf("negative value %f for ID %d", record.Value, record.ID))
		}
	}

	return errors
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
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID        int
	Name      string
	Value     float64
	Timestamp string
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
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

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
		}

		record, err := parseRow(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid data found in file")
	}

	return records, nil
}

func parseRow(row []string, lineNumber int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return record, fmt.Errorf("empty name at line %d", lineNumber)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}
	record.Value = value

	record.Timestamp = strings.TrimSpace(row[3])
	if record.Timestamp == "" {
		return record, fmt.Errorf("empty timestamp at line %d", lineNumber)
	}

	return record, nil
}

func ValidateRecords(records []DataRecord) error {
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
			return fmt.Errorf("negative value %f for ID %d", record.Value, record.ID)
		}
	}

	return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, float64) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var min, max float64

	for i, record := range records {
		sum += record.Value
		if i == 0 {
			min = record.Value
			max = record.Value
		} else {
			if record.Value < min {
				min = record.Value
			}
			if record.Value > max {
				max = record.Value
			}
		}
	}

	average := sum / float64(len(records))
	return average, min, max
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := ProcessCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	if err := ValidateRecords(records); err != nil {
		fmt.Printf("Data validation failed: %v\n", err)
		os.Exit(1)
	}

	avg, min, max := CalculateStatistics(records)
	fmt.Printf("Processed %d records successfully\n", len(records))
	fmt.Printf("Statistics - Average: %.2f, Min: %.2f, Max: %.2f\n", avg, min, max)
}