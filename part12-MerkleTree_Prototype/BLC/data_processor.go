
package main

import (
	"errors"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     string
	Timestamp time.Time
	Processed bool
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if len(record.Value) > 1000 {
		return errors.New("value exceeds maximum length")
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	record.Value = strings.ToUpper(strings.TrimSpace(record.Value))
	record.Processed = true
	return record
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, err
		}
		processed = append(processed, TransformRecord(record))
	}
	return processed, nil
}