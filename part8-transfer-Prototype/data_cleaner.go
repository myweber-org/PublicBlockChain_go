
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Name  string
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) {
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
	seen := make(map[string]bool)
	result := make([]DataRecord, 0)

	for _, record := range dc.records {
		key := fmt.Sprintf("%d|%s", record.ID, strings.ToLower(record.Email))
		if !seen[key] {
			seen[key] = true
			result = append(result, record)
		}
	}

	dc.records = result
	return result
}

func (dc *DataCleaner) ValidateEmails() []DataRecord {
	validRecords := make([]DataRecord, 0)

	for _, record := range dc.records {
		if strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".") {
			validRecords = append(validRecords, record)
		}
	}

	return validRecords
}

func (dc *DataCleaner) GetRecordCount() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	cleaner.AddRecord(DataRecord{ID: 1, Email: "user@example.com", Name: "John"})
	cleaner.AddRecord(DataRecord{ID: 2, Email: "user@example.com", Name: "John"})
	cleaner.AddRecord(DataRecord{ID: 3, Email: "invalid-email", Name: "Jane"})
	cleaner.AddRecord(DataRecord{ID: 4, Email: "another@test.org", Name: "Bob"})

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

	cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", cleaner.GetRecordCount())

	valid := cleaner.ValidateEmails()
	fmt.Printf("Valid email records: %d\n", len(valid))
}