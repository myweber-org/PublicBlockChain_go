package main

import (
	"strings"
)

// CleanString removes duplicate spaces and trims leading/trailing whitespace
func CleanString(input string) string {
	// Trim spaces from start and end
	trimmed := strings.TrimSpace(input)
	
	// Split by spaces and filter out empty strings
	words := strings.Fields(trimmed)
	
	// Join back with single spaces
	return strings.Join(words, " ")
}

// RemoveDuplicates removes duplicate entries from a slice of strings
func RemoveDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// CleanSlice applies CleanString to each element and removes duplicates
func CleanSlice(items []string) []string {
	cleaned := make([]string, len(items))
	
	for i, item := range items {
		cleaned[i] = CleanString(item)
	}
	
	return RemoveDuplicates(cleaned)
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

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
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
	deduped := DeduplicateRecords(records)
	var cleaned []DataRecord

	for _, record := range deduped {
		record.Valid = ValidateEmail(record.Email)
		cleaned = append(cleaned, record)
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
		{5, "ANOTHER@test.org", false},
	}

	cleaned := CleanData(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))

	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", record.ID, record.Email, record.Valid)
	}
}