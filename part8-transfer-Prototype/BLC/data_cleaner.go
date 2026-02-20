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

	strings := []string{"apple", "banana", "apple", "orange", "banana"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
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
		if !seen[email] && len(email) > 0 {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func ValidateRecords(records []DataRecord) ([]DataRecord, []DataRecord) {
	var valid []DataRecord
	var invalid []DataRecord
	
	for _, record := range records {
		if record.Valid && record.ID > 0 && strings.Contains(record.Email, "@") {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}
	return valid, invalid
}

func CleanData(input []string) []string {
	deduped := DeduplicateEmails(input)
	var cleaned []string
	for _, item := range deduped {
		cleaned = append(cleaned, strings.TrimSpace(item))
	}
	return cleaned
}

func main() {
	testEmails := []string{
		"user@example.com",
		"USER@example.com",
		"user@example.com",
		"",
		"  test@domain.org  ",
	}
	
	cleaned := CleanData(testEmails)
	fmt.Println("Cleaned emails:", cleaned)
	
	records := []DataRecord{
		{1, "alice@test.com", true},
		{2, "bob@test.com", false},
		{0, "invalid", true},
		{3, "charlie@test.com", true},
	}
	
	valid, invalid := ValidateRecords(records)
	fmt.Printf("Valid records: %d\n", len(valid))
	fmt.Printf("Invalid records: %d\n", len(invalid))
}