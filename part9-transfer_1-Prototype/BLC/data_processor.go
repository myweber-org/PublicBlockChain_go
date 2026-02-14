
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func ProcessCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]Record, 0)

	// Skip header
	_, err = reader.Read()
	if err != nil && err != io.EOF {
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

		if len(row) != 3 {
			return nil, errors.New("invalid row format")
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %w", err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value format: %w", err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  row[1],
			Value: value,
		})
	}

	return records, nil
}

func ValidateRecords(records []Record) error {
	if len(records) == 0 {
		return errors.New("no records to validate")
	}

	seenIDs := make(map[int]bool)
	for _, rec := range records {
		if rec.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", rec.ID)
		}
		if rec.Name == "" {
			return fmt.Errorf("record %d has empty name", rec.ID)
		}
		if rec.Value < 0 {
			return fmt.Errorf("record %d has negative value", rec.ID)
		}
		if seenIDs[rec.ID] {
			return fmt.Errorf("duplicate ID found: %d", rec.ID)
		}
		seenIDs[rec.ID] = true
	}

	return nil
}

func CalculateTotalValue(records []Record) float64 {
	var total float64
	for _, rec := range records {
		total += rec.Value
	}
	return total
}