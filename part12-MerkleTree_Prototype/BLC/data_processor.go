package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

type DataRecord struct {
	ID    string
	Name  string
	Value float64
	Valid bool
}

func ParseCSVData(input string) ([]DataRecord, error) {
	reader := csv.NewReader(strings.NewReader(input))
	records := []DataRecord{}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 3 {
			continue
		}

		record := DataRecord{
			ID:   strings.TrimSpace(row[0]),
			Name: strings.TrimSpace(row[1]),
		}

		var value float64
		if _, err := fmt.Sscanf(row[2], "%f", &value); err == nil {
			record.Value = value
			record.Valid = true
		}

		records = append(records, record)
	}

	return records, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	for _, record := range records {
		if record.Valid && record.Value > 0 {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateTotal(records []DataRecord) float64 {
	total := 0.0
	for _, record := range records {
		total += record.Value
	}
	return total
}

func main() {
	csvData := `001,ProductA,25.50
002,ProductB,18.75
003,ProductC,invalid
004,ProductD,42.00
005,ProductE,-5.00`

	records, err := ParseCSVData(csvData)
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		return
	}

	validRecords := FilterValidRecords(records)
	total := CalculateTotal(validRecords)

	fmt.Printf("Parsed %d records\n", len(records))
	fmt.Printf("Found %d valid records\n", len(validRecords))
	fmt.Printf("Total value: %.2f\n", total)
}