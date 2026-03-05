package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6, 7}
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if strings.Contains(record.Email, "@") && len(record.Email) > 3 {
			record.Valid = true
			valid = append(valid, record)
		}
	}
	return valid
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
	}

	unique := RemoveDuplicates(records)
	validated := ValidateEmails(unique)

	fmt.Printf("Processed %d records\n", len(validated))
	for _, r := range validated {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}package main

import (
    "fmt"
    "strings"
)

type DataRecord struct {
    ID    int
    Email string
    Valid bool
}

func deduplicateEmails(emails []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    for _, email := range emails {
        if !seen[email] {
            seen[email] = true
            result = append(result, email)
        }
    }
    return result
}

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func cleanData(records []DataRecord) []DataRecord {
    emailSet := make(map[string]bool)
    cleaned := []DataRecord{}
    
    for _, record := range records {
        if !emailSet[record.Email] && validateEmail(record.Email) {
            emailSet[record.Email] = true
            record.Valid = true
            cleaned = append(cleaned, record)
        }
    }
    return cleaned
}

func main() {
    records := []DataRecord{
        {1, "test@example.com", false},
        {2, "invalid-email", false},
        {3, "test@example.com", false},
        {4, "user@domain.org", false},
    }
    
    cleaned := cleanData(records)
    fmt.Printf("Cleaned %d records\n", len(cleaned))
    for _, r := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
    }
}
package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Score int
}

type DataCleaner struct {
	records map[string]DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make(map[string]DataRecord),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if !strings.Contains(record.Email, "@") {
		return errors.New("invalid email format")
	}
	if record.Score < 0 || record.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}

	if _, exists := dc.records[record.ID]; exists {
		return fmt.Errorf("duplicate record ID: %s", record.ID)
	}

	dc.records[record.ID] = record
	return nil
}

func (dc *DataCleaner) RemoveRecord(id string) bool {
	if _, exists := dc.records[id]; exists {
		delete(dc.records, id)
		return true
	}
	return false
}

func (dc *DataCleaner) GetValidRecords() []DataRecord {
	var validRecords []DataRecord
	for _, record := range dc.records {
		validRecords = append(validRecords, record)
	}
	return validRecords
}

func (dc *DataCleaner) CountRecords() int {
	return len(dc.records)
}

func (dc *DataCleaner) FindByEmailDomain(domain string) []DataRecord {
	var results []DataRecord
	for _, record := range dc.records {
		if strings.HasSuffix(record.Email, domain) {
			results = append(results, record)
		}
	}
	return results
}

func main() {
	cleaner := NewDataCleaner()

	sampleRecords := []DataRecord{
		{ID: "001", Email: "user1@example.com", Score: 85},
		{ID: "002", Email: "user2@test.org", Score: 92},
		{ID: "003", Email: "user3@example.com", Score: 78},
	}

	for _, record := range sampleRecords {
		if err := cleaner.AddRecord(record); err != nil {
			fmt.Printf("Failed to add record %s: %v\n", record.ID, err)
		}
	}

	fmt.Printf("Total valid records: %d\n", cleaner.CountRecords())

	exampleUsers := cleaner.FindByEmailDomain("@example.com")
	fmt.Printf("Users with example.com domain: %d\n", len(exampleUsers))

	cleaner.RemoveRecord("002")
	fmt.Printf("Records after removal: %d\n", cleaner.CountRecords())
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmails(records []DataRecord) []DataRecord {
	for i := range records {
		email := records[i].Email
		records[i].Valid = strings.Contains(email, "@") && strings.Contains(email, ".")
	}
	return records
}

func processData(records []DataRecord) []DataRecord {
	records = deduplicateRecords(records)
	records = validateEmails(records)
	return records
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
	}

	processed := processData(sampleData)

	for _, record := range processed {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
	}
}