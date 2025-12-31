package main

import (
	"errors"
	"strings"
)

type DataRecord struct {
	ID    string
	Value string
	Valid bool
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if len(record.Value) > 100 {
		return errors.New("value exceeds maximum length")
	}
	return nil
}

func TransformValue(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, rec := range records {
		if err := ValidateRecord(rec); err != nil {
			return nil, err
		}
		rec.Value = TransformValue(rec.Value)
		rec.Valid = true
		processed = append(processed, rec)
	}
	return processed, nil
}