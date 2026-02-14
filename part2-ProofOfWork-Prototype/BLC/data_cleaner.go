
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

func deduplicateEmails(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			record.Email = email
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	for _, record := range records {
		if validateEmail(record.Email) {
			record.Valid = true
			validRecords = append(validRecords, record)
		}
	}
	return deduplicateEmails(validRecords)
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "test@domain.org", false},
		{4, "invalid-email", false},
		{5, "test@domain.org", false},
	}

	cleaned := processRecords(sampleData)
	fmt.Printf("Processed %d records, %d valid unique records found\n", len(sampleData), len(cleaned))
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", record.ID, record.Email, record.Valid)
	}
}package datautils

func RemoveDuplicates(input []int) []int {
    if len(input) == 0 {
        return input
    }
    
    seen := make(map[int]bool)
    result := make([]int, 0, len(input))
    
    for _, value := range input {
        if !seen[value] {
            seen[value] = true
            result = append(result, value)
        }
    }
    
    return result
}