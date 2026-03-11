
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)

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
			return nil, errors.New("invalid row format")
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %w", err)
		}

		name := strings.TrimSpace(row[1])

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value format: %w", err)
		}

		valid := strings.ToLower(strings.TrimSpace(row[3])) == "true"

		record := DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
			Valid: valid,
		}
		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, error) {
	validRecords := make([]DataRecord, 0)
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return nil, fmt.Errorf("invalid ID %d: must be positive", record.ID)
		}

		if seenIDs[record.ID] {
			return nil, fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Name == "" {
			return nil, errors.New("name cannot be empty")
		}

		if record.Value < 0 {
			return nil, fmt.Errorf("negative value not allowed: %f", record.Value)
		}

		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}

	return validRecords, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
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

func ProcessDataFile(filename string) error {
	records, err := ParseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	validRecords, err := ValidateRecords(records)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	avg, max, count := CalculateStatistics(validRecords)
	fmt.Printf("Processing complete:\n")
	fmt.Printf("  Total records: %d\n", len(records))
	fmt.Printf("  Valid records: %d\n", count)
	fmt.Printf("  Average value: %.2f\n", avg)
	fmt.Printf("  Maximum value: %.2f\n", max)

	return nil
}