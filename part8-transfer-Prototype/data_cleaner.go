
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
package main

import (
	"errors"
	"fmt"
	"strings"
)

type Record struct {
	ID    int
	Email string
	Valid bool
}

func DeduplicateEmails(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, rec := range records {
		email := strings.ToLower(strings.TrimSpace(rec.Email))
		if !seen[email] && email != "" {
			seen[email] = true
			unique = append(unique, Record{
				ID:    rec.ID,
				Email: email,
				Valid: rec.Valid,
			})
		}
	}
	return unique
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("email must contain @ symbol")
	}
	if !strings.Contains(email, ".") {
		return errors.New("email must contain domain")
	}
	return nil
}

func CleanRecords(records []Record) ([]Record, error) {
	cleaned := DeduplicateEmails(records)
	for i := range cleaned {
		if err := ValidateEmail(cleaned[i].Email); err != nil {
			cleaned[i].Valid = false
			fmt.Printf("Warning: Record ID %d invalid: %v\n", cleaned[i].ID, err)
		} else {
			cleaned[i].Valid = true
		}
	}
	return cleaned, nil
}

func main() {
	sampleData := []Record{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
		{5, "", false},
	}

	cleaned, err := CleanRecords(sampleData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Cleaned %d records:\n", len(cleaned))
	for _, rec := range cleaned {
		status := "valid"
		if !rec.Valid {
			status = "invalid"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", rec.ID, rec.Email, status)
	}
}