
package data_processor

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

type DataRecord struct {
	ID        string  `json:"id"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
	Category  string  `json:"category"`
}

func ParseAndValidateJSON(rawData []byte) (*DataRecord, error) {
	var record DataRecord
	
	if err := json.Unmarshal(rawData, &record); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var validationErrors []ValidationError

	if strings.TrimSpace(record.ID) == "" {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "id",
			Message: "cannot be empty",
		})
	}

	if record.Value < 0 {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "value",
			Message: "must be non-negative",
		})
	}

	if record.Timestamp <= 0 {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "timestamp",
			Message: "must be positive integer",
		})
	}

	if !isValidCategory(record.Category) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "category",
			Message: "invalid category specified",
		})
	}

	if len(validationErrors) > 0 {
		var errorMessages []string
		for _, err := range validationErrors {
			errorMessages = append(errorMessages, err.Error())
		}
		return nil, fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
	}

	return &record, nil
}

func isValidCategory(category string) bool {
	validCategories := map[string]bool{
		"standard": true,
		"premium":  true,
		"legacy":   true,
	}
	return validCategories[category]
}
package main

import (
	"strings"
	"unicode"
)

func ProcessInput(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	trimmed := strings.TrimSpace(input)
	var cleaned strings.Builder

	for _, r := range trimmed {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			cleaned.WriteRune(r)
		}
	}

	result := strings.Join(strings.Fields(cleaned.String()), " ")
	return result, nil
}
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

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
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
    total := 0.0
    for _, rec := range records {
        total += rec.Value
    }
    return total
}