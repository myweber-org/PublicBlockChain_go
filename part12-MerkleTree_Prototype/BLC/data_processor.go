
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
package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataRecord
	for i, row := range records {
		if len(row) < 4 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(row[1])

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			continue
		}

		active := strings.ToLower(strings.TrimSpace(row[3])) == "true"

		data = append(data, DataRecord{
			ID:     id,
			Name:   name,
			Value:  value,
			Active: active,
		})
	}

	return data, nil
}

func ValidateData(records []DataRecord) error {
	if len(records) == 0 {
		return errors.New("no data records found")
	}

	seenIDs := make(map[int]bool)
	for _, record := range records {
		if seenIDs[record.ID] {
			return errors.New("duplicate ID found: " + strconv.Itoa(record.ID))
		}
		seenIDs[record.ID] = true

		if record.ID <= 0 {
			return errors.New("invalid ID: must be positive")
		}

		if strings.TrimSpace(record.Name) == "" {
			return errors.New("empty name field")
		}

		if record.Value < 0 {
			return errors.New("negative value not allowed")
		}
	}

	return nil
}

func FilterActiveRecords(records []DataRecord) []DataRecord {
	var active []DataRecord
	for _, record := range records {
		if record.Active {
			active = append(active, record)
		}
	}
	return active
}

func CalculateTotalValue(records []DataRecord) float64 {
	var total float64
	for _, record := range records {
		total += record.Value
	}
	return total
}package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	allowedPattern *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9\s\-_@.]+$`)
	return &DataProcessor{allowedPattern: pattern}
}

func (dp *DataProcessor) SanitizeInput(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	if !dp.allowedPattern.MatchString(trimmed) {
		return "", false
	}

	return trimmed, true
}

func (dp *DataProcessor) ProcessUserData(rawData []string) []string {
	var cleanData []string
	for _, item := range rawData {
		if sanitized, ok := dp.SanitizeInput(item); ok {
			cleanData = append(cleanData, sanitized)
		}
	}
	return cleanData
}