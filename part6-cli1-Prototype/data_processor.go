
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

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}

		if len(row) < 3 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) []DataRecord {
	validRecords := make([]DataRecord, 0)
	for _, record := range records {
		if record.ID > 0 && record.Name != "" && record.Value >= 0 {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}

	total := 0.0
	for _, record := range records {
		total += record.Value
	}
	return total / float64(len(records))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		return
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	validRecords := ValidateRecords(records)
	average := CalculateAverage(validRecords)

	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Valid records: %d\n", len(validRecords))
	fmt.Printf("Average value: %.2f\n", average)
}