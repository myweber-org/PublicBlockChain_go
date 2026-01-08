
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(normalized) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && !strings.ContainsAny(item, "!@#$%")
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	data := []string{"apple", "Apple", "banana", "", "cherry!", "banana", "date "}
	
	cleaned := cleaner.RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
	
	cleaner.Reset()
	testData := []string{"test1", "test2", "test1"}
	fmt.Printf("Reset test: %v\n", cleaner.RemoveDuplicates(testData))
}
package main

import (
    "fmt"
    "strings"
)

type DataRecord struct {
    ID    int
    Email string
    Phone string
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
    seen := make(map[string]bool)
    unique := []DataRecord{}
    
    for _, record := range records {
        key := fmt.Sprintf("%s|%s", record.Email, record.Phone)
        if !seen[key] {
            seen[key] = true
            unique = append(unique, record)
        }
    }
    return unique
}

func ValidateEmail(email string) bool {
    if len(email) == 0 {
        return false
    }
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func ValidatePhone(phone string) bool {
    if len(phone) == 0 {
        return false
    }
    for _, ch := range phone {
        if ch < '0' || ch > '9' {
            return false
        }
    }
    return len(phone) >= 10
}

func CleanData(records []DataRecord) []DataRecord {
    validRecords := []DataRecord{}
    
    for _, record := range records {
        if ValidateEmail(record.Email) && ValidatePhone(record.Phone) {
            validRecords = append(validRecords, record)
        }
    }
    
    return DeduplicateRecords(validRecords)
}

func main() {
    sampleData := []DataRecord{
        {1, "test@example.com", "1234567890"},
        {2, "test@example.com", "1234567890"},
        {3, "invalid-email", "1234567890"},
        {4, "another@test.com", "not-a-phone"},
        {5, "valid@email.org", "0987654321"},
    }
    
    cleaned := CleanData(sampleData)
    
    fmt.Printf("Original records: %d\n", len(sampleData))
    fmt.Printf("Cleaned records: %d\n", len(cleaned))
    
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", 
            record.ID, record.Email, record.Phone)
    }
}
package main

import "fmt"

func removeDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "cherry", "banana", "date"}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
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
	ID    int
	Name  string
	Email string
	Score float64
}

func cleanCSVData(inputPath string, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
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

	headers = append(headers, "Valid")
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	lineNumber := 1
	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Warning: line %d: %v\n", lineNumber, err)
			continue
		}

		record, validationErr := validateRecord(row)
		isValid := validationErr == nil

		outputRow := make([]string, len(row)+1)
		copy(outputRow, row)
		outputRow[len(row)] = strconv.FormatBool(isValid)

		if isValid {
			fmt.Printf("Processed record ID %d: %s\n", record.ID, record.Name)
		} else {
			fmt.Printf("Invalid record at line %d: %v\n", lineNumber, validationErr)
		}

		if err := writer.Write(outputRow); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

func validateRecord(row []string) (Record, error) {
	if len(row) < 4 {
		return Record{}, fmt.Errorf("insufficient columns")
	}

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID: %w", err)
	}

	name := strings.TrimSpace(row[1])
	if name == "" {
		return Record{}, fmt.Errorf("name cannot be empty")
	}

	email := strings.TrimSpace(row[2])
	if !strings.Contains(email, "@") {
		return Record{}, fmt.Errorf("invalid email format")
	}

	score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid score: %w", err)
	}

	if score < 0 || score > 100 {
		return Record{}, fmt.Errorf("score out of range (0-100)")
	}

	return Record{
		ID:    id,
		Name:  name,
		Email: email,
		Score: score,
	}, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSVData(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data cleaning completed successfully")
}