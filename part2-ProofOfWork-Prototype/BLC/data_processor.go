package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ProcessRecord(record DataRecord) (string, error) {
	if record.ID <= 0 {
		return "", errors.New("invalid record ID")
	}

	if strings.TrimSpace(record.Name) == "" {
		return "", errors.New("record name cannot be empty")
	}

	if record.Value < 0 {
		return "", errors.New("record value cannot be negative")
	}

	processedName := strings.ToUpper(record.Name)
	formattedValue := fmt.Sprintf("%.2f", record.Value)

	result := fmt.Sprintf("Processed: ID=%d, NAME=%s, VALUE=%s",
		record.ID, processedName, formattedValue)

	return result, nil
}

func ValidateAndProcess(records []DataRecord) ([]string, []error) {
	var results []string
	var errs []error

	for _, record := range records {
		result, err := ProcessRecord(record)
		if err != nil {
			errs = append(errs, fmt.Errorf("record %d: %w", record.ID, err))
			continue
		}
		results = append(results, result)
	}

	return results, errs
}

func main() {
	records := []DataRecord{
		{ID: 1, Name: "record_one", Value: 100.50},
		{ID: 2, Name: "record_two", Value: -5.0},
		{ID: 0, Name: "record_three", Value: 75.25},
		{ID: 4, Name: "", Value: 200.0},
		{ID: 5, Name: "record_five", Value: 300.75},
	}

	results, errs := ValidateAndProcess(records)

	fmt.Println("Processing Results:")
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\nErrors:")
	for _, err := range errs {
		fmt.Println(err)
	}
}