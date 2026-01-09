package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	Username string
	Email    string
	Age      int
}

func ValidateProfile(profile UserProfile) error {
	if strings.TrimSpace(profile.Username) == "" {
		return errors.New("username cannot be empty")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if profile.Age < 0 || profile.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func NormalizeProfile(profile UserProfile) UserProfile {
	normalized := profile
	normalized.Username = strings.ToLower(strings.TrimSpace(profile.Username))
	normalized.Email = strings.ToLower(strings.TrimSpace(profile.Email))
	return normalized
}

func ProcessUserData(profile UserProfile) (UserProfile, error) {
	if err := ValidateProfile(profile); err != nil {
		return UserProfile{}, err
	}
	return NormalizeProfile(profile), nil
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

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)

	lineNumber := 0
	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		record := DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		}

		if err := validateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed at line %d: %w", lineNumber, err)
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		return nil, errors.New("no valid records found in file")
	}

	return records, nil
}

func validateRecord(record DataRecord) error {
	if record.ID <= 0 {
		return errors.New("ID must be positive")
	}
	if record.Value < 0 {
		return errors.New("value cannot be negative")
	}
	return nil
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

func FilterRecords(records []DataRecord, minValue float64) []DataRecord {
	filtered := make([]DataRecord, 0)
	for _, record := range records {
		if record.Value >= minValue {
			filtered = append(filtered, record)
		}
	}
	return filtered
}