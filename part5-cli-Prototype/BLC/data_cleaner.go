
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

func DeduplicateEmails(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] && email != "" {
			seen[email] = true
			record.Email = email
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmailFormat(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func CleanDataset(records []DataRecord) []DataRecord {
	records = DeduplicateEmails(records)
	var cleaned []DataRecord
	for _, record := range records {
		record.Valid = ValidateEmailFormat(record.Email)
		if record.Valid {
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	dataset := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain", false},
		{5, "user@example.com", false},
		{6, "new@test.org", false},
	}

	cleaned := CleanDataset(dataset)
	fmt.Printf("Original: %d records\n", len(dataset))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s\n", record.ID, record.Email)
	}
}