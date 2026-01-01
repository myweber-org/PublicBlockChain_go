package data

import (
	"errors"
	"strings"
	"time"
)

type Record struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(r Record) error {
	if r.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if r.Value < 0 {
		return errors.New("record value must be non-negative")
	}
	if r.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(r Record, multiplier float64) (Record, error) {
	if err := ValidateRecord(r); err != nil {
		return Record{}, err
	}

	transformed := r
	transformed.Value = r.Value * multiplier

	for i, tag := range transformed.Tags {
		transformed.Tags[i] = strings.ToUpper(strings.TrimSpace(tag))
	}

	return transformed, nil
}

func FilterRecords(records []Record, minValue float64) []Record {
	var filtered []Record
	for _, r := range records {
		if r.Value >= minValue {
			filtered = append(filtered, r)
		}
	}
	return filtered
}