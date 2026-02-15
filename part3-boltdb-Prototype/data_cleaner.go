
package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
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
	validRecords := make([]DataRecord, 0, len(dc.records))
	for _, record := range dc.records {
		validRecords = append(validRecords, record)
	}
	return validRecords
}

func (dc *DataCleaner) Count() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	sampleRecords := []DataRecord{
		{ID: "001", Email: "user1@example.com", Score: 85},
		{ID: "002", Email: "user2@example.com", Score: 92},
		{ID: "003", Email: "invalid-email", Score: 75},
		{ID: "001", Email: "duplicate@example.com", Score: 88},
		{ID: "004", Email: "user4@example.com", Score: 105},
	}

	for _, record := range sampleRecords {
		err := cleaner.AddRecord(record)
		if err != nil {
			fmt.Printf("Failed to add record %s: %v\n", record.ID, err)
		} else {
			fmt.Printf("Added record %s successfully\n", record.ID)
		}
	}

	fmt.Printf("\nTotal valid records: %d\n", cleaner.Count())
	fmt.Println("Valid records:")
	for _, record := range cleaner.GetValidRecords() {
		fmt.Printf("  ID: %s, Email: %s, Score: %d\n", record.ID, record.Email, record.Score)
	}
}
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

func deduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
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
	var cleaned []DataRecord

	for _, record := range records {
		record.Email = strings.ToLower(strings.TrimSpace(record.Email))
		if validateEmail(record.Email) && !emailSet[record.Email] {
			record.Valid = true
			emailSet[record.Email] = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
	}

	cleaned := cleanData(records)
	fmt.Printf("Cleaned %d records, kept %d\n", len(records), len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}
package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}