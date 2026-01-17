
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
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Value = record.Value * 1.1
	transformed.Tags = append(record.Tags, "processed")
	return transformed
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

func GenerateSummary(records []DataRecord) string {
	if len(records) == 0 {
		return "No records to summarize"
	}
	
	var total float64
	var tagCount int
	for _, record := range records {
		total += record.Value
		tagCount += len(record.Tags)
	}
	
	avg := total / float64(len(records))
	return fmt.Sprintf("Processed %d records. Average value: %.2f, Total tags: %d", 
		len(records), avg, tagCount)
}

func FilterByTag(records []DataRecord, tag string) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		for _, t := range record.Tags {
			if strings.EqualFold(t, tag) {
				filtered = append(filtered, record)
				break
			}
		}
	}
	return filtered
}