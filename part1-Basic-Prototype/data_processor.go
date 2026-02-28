
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
			Email: strings.TrimSpace(row[2]),
			Valid: validateRecord(row[0], row[1], row[2]),
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(id, name, email string) bool {
	if id == "" || name == "" || email == "" {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	return true
}

func GenerateReport(records []DataRecord) {
	validCount := 0
	invalidCount := 0

	fmt.Println("=== DATA PROCESSING REPORT ===")
	for _, record := range records {
		if record.Valid {
			validCount++
			fmt.Printf("✓ Valid: %s - %s\n", record.ID, record.Name)
		} else {
			invalidCount++
			fmt.Printf("✗ Invalid: %s - %s (%s)\n", record.ID, record.Name, record.Email)
		}
	}

	fmt.Printf("\nSummary: %d valid, %d invalid, %d total\n", validCount, invalidCount, len(records))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	GenerateReport(records)
}