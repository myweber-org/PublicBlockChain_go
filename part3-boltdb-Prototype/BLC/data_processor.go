
package main

import (
    "fmt"
)

// FilterAndDouble filters out even numbers from the input slice,
// doubles the remaining odd numbers, and returns the new slice.
func FilterAndDouble(numbers []int) []int {
    var result []int
    for _, num := range numbers {
        if num%2 != 0 {
            result = append(result, num*2)
        }
    }
    return result
}

func main() {
    input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    output := FilterAndDouble(input)
    fmt.Println("Original:", input)
    fmt.Println("Processed:", output)
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
	Category  string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value cannot be negative")
	}
	if record.Category == "" {
		return errors.New("category cannot be empty")
	}
	if record.Timestamp.After(time.Now()) {
		return errors.New("timestamp cannot be in the future")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	return DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * 1.1,
		Timestamp: record.Timestamp.UTC(),
		Category:  strings.ToLower(record.Category),
	}
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
		{
			ID:        "abc123",
			Value:     100.0,
			Timestamp: time.Now().Add(-time.Hour),
			Category:  "SAMPLE",
		},
		{
			ID:        "def456",
			Value:     200.0,
			Timestamp: time.Now().Add(-2 * time.Hour),
			Category:  "TEST",
		},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	fmt.Printf("Processed %d records\n", len(processed))
	for _, record := range processed {
		fmt.Printf("ID: %s, Value: %.2f, Category: %s\n", 
			record.ID, record.Value, record.Category)
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

type Record struct {
    ID      int
    Name    string
    Value   float64
    Active  bool
}

func ProcessCSV(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
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

        active, err := strconv.ParseBool(row[3])
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

func CalculateStats(records []Record) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var activeCount int
    var maxValue float64

    for _, record := range records {
        sum += record.Value
        if record.Value > maxValue {
            maxValue = record.Value
        }
        if record.Active {
            activeCount++
        }
    }

    average := sum / float64(len(records))
    return average, maxValue, activeCount
}

func FilterRecords(records []Record, minValue float64) []Record {
    var filtered []Record
    for _, record := range records {
        if record.Value >= minValue {
            filtered = append(filtered, record)
        }
    }
    return filtered
}
package main

import (
	"errors"
	"strings"
	"unicode"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateAndNormalizeUserData(data UserData) (UserData, error) {
	var normalized UserData

	// Normalize username
	normalized.Username = strings.TrimSpace(data.Username)
	if len(normalized.Username) < 3 {
		return UserData{}, errors.New("username must be at least 3 characters")
	}
	for _, r := range normalized.Username {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' && r != '-' {
			return UserData{}, errors.New("username contains invalid characters")
		}
	}

	// Normalize email
	normalized.Email = strings.ToLower(strings.TrimSpace(data.Email))
	if !strings.Contains(normalized.Email, "@") {
		return UserData{}, errors.New("invalid email format")
	}

	// Validate age
	if data.Age < 0 || data.Age > 150 {
		return UserData{}, errors.New("age must be between 0 and 150")
	}
	normalized.Age = data.Age

	return normalized, nil
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
	rawData := UserData{
		Username: username,
		Email:    email,
		Age:      age,
	}
	return ValidateAndNormalizeUserData(rawData)
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func parseCSVFile(filename string) ([]DataRecord, error) {
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

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
	}
	record.Active = active

	return record, nil
}

func validateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("negative value for ID %d: %f", record.ID, record.Value)
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

	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Total value: %.2f\n", totalValue)
	fmt.Printf("Active records: %d\n", activeCount)
	fmt.Printf("Average value: %.2f\n", totalValue/float64(len(records)))

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