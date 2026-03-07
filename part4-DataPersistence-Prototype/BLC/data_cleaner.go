package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
	encountered := map[int]bool{}
	result := []DataRecord{}

	for _, record := range records {
		if !encountered[record.ID] {
			encountered[record.ID] = true
			result = append(result, record)
		}
	}
	return result
}

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func CleanData(records []DataRecord) ([]DataRecord, error) {
	cleaned := RemoveDuplicates(records)
	for i, record := range cleaned {
		if err := ValidateEmail(record.Email); err != nil {
			cleaned[i].Valid = false
			fmt.Printf("Warning: Record ID %d has invalid email: %v\n", record.ID, err)
		} else {
			cleaned[i].Valid = true
		}
	}
	return cleaned, nil
}

func main() {
	sampleData := []DataRecord{
		{ID: 1, Email: "test@example.com"},
		{ID: 2, Email: "invalid-email"},
		{ID: 1, Email: "test@example.com"},
		{ID: 3, Email: "another@domain.org"},
	}

	cleaned, err := CleanData(sampleData)
	if err != nil {
		fmt.Printf("Error cleaning data: %v\n", err)
		return
	}

	fmt.Printf("Cleaned %d records:\n", len(cleaned))
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
	}
}