
package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Age   int
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) error {
	if record.Name == "" {
		return errors.New("name cannot be empty")
	}
	if record.Age < 0 || record.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	if !strings.Contains(record.Email, "@") {
		return errors.New("invalid email format")
	}

	dc.records = append(dc.records, record)
	return nil
}

func (dc *DataCleaner) RemoveDuplicates() {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range dc.records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}

	dc.records = unique
}

func (dc *DataCleaner) ValidateAll() []error {
	var errs []error

	for i, record := range dc.records {
		if record.Name == "" {
			errs = append(errs, fmt.Errorf("record %d: name is empty", i))
		}
		if record.Age < 0 {
			errs = append(errs, fmt.Errorf("record %d: age is negative", i))
		}
		if !strings.Contains(record.Email, "@") {
			errs = append(errs, fmt.Errorf("record %d: invalid email", i))
		}
	}

	return errs
}

func (dc *DataCleaner) GetRecords() []DataRecord {
	return dc.records
}

func (dc *DataCleaner) Count() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", 30},
		{2, "Jane Smith", "jane@example.com", 25},
		{3, "John Doe", "john@example.com", 30},
		{4, "Invalid User", "invalid-email", -5},
	}

	for _, record := range sampleData {
		if err := cleaner.AddRecord(record); err != nil {
			fmt.Printf("Failed to add record: %v\n", err)
		}
	}

	fmt.Printf("Initial count: %d\n", cleaner.Count())

	cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", cleaner.Count())

	validationErrors := cleaner.ValidateAll()
	if len(validationErrors) > 0 {
		fmt.Println("Validation errors found:")
		for _, err := range validationErrors {
			fmt.Printf("  - %v\n", err)
		}
	}

	fmt.Println("Final records:")
	for _, record := range cleaner.GetRecords() {
		fmt.Printf("  ID: %d, Name: %s, Email: %s, Age: %d\n",
			record.ID, record.Name, record.Email, record.Age)
	}
}