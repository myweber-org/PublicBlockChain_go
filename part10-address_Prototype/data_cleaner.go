
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
		normalized := strings.ToLower(strings.TrimSpace(email))
		if !seen[normalized] {
			seen[normalized] = true
			result = append(result, normalized)
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
	return strings.Contains(parts[1], ".")
}

func CleanData(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	cleaned := []DataRecord{}
	for _, record := range records {
		record.Email = strings.ToLower(strings.TrimSpace(record.Email))
		if ValidateEmail(record.Email) && !emailSet[record.Email] {
			record.Valid = true
			emailSet[record.Email] = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	emails := []string{"test@example.com", "TEST@example.com", "invalid", "another@test.org"}
	unique := DeduplicateEmails(emails)
	fmt.Println("Deduplicated emails:", unique)

	records := []DataRecord{
		{1, "user@domain.com", false},
		{2, "USER@domain.com", false},
		{3, "bad-email", false},
	}
	cleaned := CleanData(records)
	fmt.Printf("Cleaned records: %+v\n", cleaned)
}