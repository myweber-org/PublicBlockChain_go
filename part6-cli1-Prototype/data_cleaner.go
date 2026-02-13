
package main

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
	input := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(input)
	fmt.Printf("Original: %v\n", input)
	fmt.Printf("Cleaned: %v\n", cleaned)
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
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	dc.seen = make(map[string]bool)
	var result []string
	for _, item := range items {
		if !dc.IsDuplicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"Apple", "apple", " BANANA ", "banana", "Cherry", "cherry "}
	
	fmt.Println("Original data:", data)
	
	uniqueData := cleaner.RemoveDuplicates(data)
	fmt.Println("Deduplicated data:", uniqueData)
	
	cleaner.Reset()
	
	testValue := "  TEST  "
	normalized := cleaner.Normalize(testValue)
	fmt.Printf("Normalized '%s' to '%s'\n", testValue, normalized)
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

func DeduplicateEmails(emails []string) []string {
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

func ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return true
}

func CleanData(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	var cleaned []DataRecord

	for _, record := range records {
		record.Email = strings.ToLower(strings.TrimSpace(record.Email))
		record.Valid = ValidateEmail(record.Email)

		if record.Valid && !emailSet[record.Email] {
			emailSet[record.Email] = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "  ADMIN@EXAMPLE.COM  ", false},
		{3, "invalid-email", false},
		{4, "user@example.com", false},
		{5, "test@domain.org", false},
	}

	cleaned := CleanData(records)
	fmt.Printf("Original: %d records\n", len(records))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))

	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedRecords int
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{processedRecords: 0}
}

func (dc *DataCleaner) RemoveDuplicates(data []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range data {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
			dc.processedRecords++
		}
	}
	return result
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if parts[0] == "" || parts[1] == "" {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

func (dc *DataCleaner) GetStats() string {
	return fmt.Sprintf("Processed records: %d", dc.processedRecords)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{
		"  john@example.com  ",
		"john@example.com",
		"sarah@test.org",
		"invalid-email",
		"",
		"  admin@company.co  ",
	}
	
	uniqueData := cleaner.RemoveDuplicates(data)
	fmt.Println("Unique items:", uniqueData)
	
	for _, email := range uniqueData {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("Valid email: %s\n", email)
		} else {
			fmt.Printf("Invalid email: %s\n", email)
		}
	}
	
	fmt.Println(cleaner.GetStats())
}