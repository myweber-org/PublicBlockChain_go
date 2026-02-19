
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
	Email   string
	Active  string
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
	headerSkipped := false

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		if !headerSkipped {
			headerSkipped = true
			continue
		}

		if len(row) < 4 {
			continue
		}

		record := DataRecord{
			ID:     strings.TrimSpace(row[0]),
			Name:   strings.TrimSpace(row[1]),
			Email:  strings.TrimSpace(row[2]),
			Active: strings.TrimSpace(row[3]),
		}

		if isValidRecord(record) {
			records = append(records, record)
		}
	}

	return records, nil
}

func isValidRecord(record DataRecord) bool {
	if record.ID == "" || record.Name == "" {
		return false
	}
	if !strings.Contains(record.Email, "@") {
		return false
	}
	return record.Active == "true" || record.Active == "false"
}

func FilterActiveRecords(records []DataRecord) []DataRecord {
	var active []DataRecord
	for _, record := range records {
		if record.Active == "true" {
			active = append(active, record)
		}
	}
	return active
}

func GenerateReport(records []DataRecord) {
	fmt.Printf("Total records processed: %d\n", len(records))
	active := FilterActiveRecords(records)
	fmt.Printf("Active records: %d\n", len(active))
	fmt.Printf("Inactive records: %d\n", len(records)-len(active))
}