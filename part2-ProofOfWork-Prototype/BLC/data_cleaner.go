
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(id int, name, email string) {
	record := DataRecord{
		ID:    id,
		Name:  strings.TrimSpace(name),
		Email: strings.TrimSpace(email),
		Valid: true,
	}
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) ValidateEmails() {
	for i := range dc.records {
		if !strings.Contains(dc.records[i].Email, "@") {
			dc.records[i].Valid = false
		}
	}
}

func (dc *DataCleaner) RemoveDuplicates() {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range dc.records {
		key := fmt.Sprintf("%d|%s|%s", record.ID, record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	dc.records = unique
}

func (dc *DataCleaner) GetValidRecords() []DataRecord {
	var valid []DataRecord
	for _, record := range dc.records {
		if record.Valid {
			valid = append(valid, record)
		}
	}
	return valid
}

func (dc *DataCleaner) PrintSummary() {
	fmt.Printf("Total records: %d\n", len(dc.records))
	valid := dc.GetValidRecords()
	fmt.Printf("Valid records: %d\n", len(valid))
}

func main() {
	cleaner := NewDataCleaner()

	cleaner.AddRecord(1, "John Doe", "john@example.com")
	cleaner.AddRecord(2, "Jane Smith", "jane@example.com")
	cleaner.AddRecord(3, "Bob Wilson", "invalid-email")
	cleaner.AddRecord(4, "John Doe", "john@example.com")

	fmt.Println("Before cleaning:")
	cleaner.PrintSummary()

	cleaner.ValidateEmails()
	cleaner.RemoveDuplicates()

	fmt.Println("\nAfter cleaning:")
	cleaner.PrintSummary()

	fmt.Println("\nValid records:")
	for _, record := range cleaner.GetValidRecords() {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}package datautils

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	cleanedHeaders := make([]string, len(headers))
	for i, h := range headers {
		cleanedHeaders[i] = strings.TrimSpace(strings.ToLower(h))
	}
	if err := writer.Write(cleanedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedField := strings.TrimSpace(field)
			if cleanedField == "" {
				cleanedField = "N/A"
			}
			cleanedRecord[i] = cleanedField
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSV(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully cleaned data. Output saved to %s\n", outputFile)
}