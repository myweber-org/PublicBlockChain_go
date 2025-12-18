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
	for _, r := range records {
		if r.IsValid {
			valid = append(valid, r)
		}
	}
	return valid
}

func GenerateSummary(records []DataRecord) {
	validCount := 0
	for _, r := range records {
		if r.IsValid {
			validCount++
		}
	}
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Valid records: %d\n", validCount)
	fmt.Printf("Invalid records: %d\n", len(records)-validCount)
}