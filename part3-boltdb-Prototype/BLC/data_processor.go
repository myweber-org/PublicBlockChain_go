
package main

import (
    "fmt"
)

// FilterAndDouble filters out even numbers from the input slice,
// doubles the remaining odd numbers, and returns the new slice.
func FilterAndDouble(numbers []int) []int {
    var result []int
    for _, num := range numbers {
        if num%2 != 0 {
            result = append(result, num*2)
        }
    }
    return result
}

func main() {
    input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    output := FilterAndDouble(input)
    fmt.Println("Original:", input)
    fmt.Println("Processed:", output)
}
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Category  string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value cannot be negative")
	}
	if record.Category == "" {
		return errors.New("category cannot be empty")
	}
	if record.Timestamp.After(time.Now()) {
		return errors.New("timestamp cannot be in the future")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	return DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * 1.1,
		Timestamp: record.Timestamp.UTC(),
		Category:  strings.ToLower(record.Category),
	}
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record))
	}
	return processed, nil
}

func main() {
	records := []DataRecord{
		{
			ID:        "abc123",
			Value:     100.0,
			Timestamp: time.Now().Add(-time.Hour),
			Category:  "SAMPLE",
		},
		{
			ID:        "def456",
			Value:     200.0,
			Timestamp: time.Now().Add(-2 * time.Hour),
			Category:  "TEST",
		},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	fmt.Printf("Processed %d records\n", len(processed))
	for _, record := range processed {
		fmt.Printf("ID: %s, Value: %.2f, Category: %s\n", 
			record.ID, record.Value, record.Category)
	}
}