
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(id int, name, email string) {
	record := DataRecord{
		ID:    id,
		Name:  strings.TrimSpace(name),
		Email: strings.TrimSpace(email),
		Valid: true,
	}
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) ValidateEmails() {
	for i := range dc.records {
		if !strings.Contains(dc.records[i].Email, "@") {
			dc.records[i].Valid = false
		}
	}
}

func (dc *DataCleaner) RemoveDuplicates() {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range dc.records {
		key := fmt.Sprintf("%d|%s|%s", record.ID, record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	dc.records = unique
}

func (dc *DataCleaner) GetValidRecords() []DataRecord {
	var valid []DataRecord
	for _, record := range dc.records {
		if record.Valid {
			valid = append(valid, record)
		}
	}
	return valid
}

func (dc *DataCleaner) PrintSummary() {
	fmt.Printf("Total records: %d\n", len(dc.records))
	valid := dc.GetValidRecords()
	fmt.Printf("Valid records: %d\n", len(valid))
}

func main() {
	cleaner := NewDataCleaner()

	cleaner.AddRecord(1, "John Doe", "john@example.com")
	cleaner.AddRecord(2, "Jane Smith", "jane@example.com")
	cleaner.AddRecord(3, "Bob Wilson", "invalid-email")
	cleaner.AddRecord(4, "John Doe", "john@example.com")

	fmt.Println("Before cleaning:")
	cleaner.PrintSummary()

	cleaner.ValidateEmails()
	cleaner.RemoveDuplicates()

	fmt.Println("\nAfter cleaning:")
	cleaner.PrintSummary()

	fmt.Println("\nValid records:")
	for _, record := range cleaner.GetValidRecords() {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}