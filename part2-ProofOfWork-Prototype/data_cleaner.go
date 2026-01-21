
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

func processRecords(records []DataRecord) []DataRecord {
	emailMap := make(map[string]bool)
	var cleanRecords []DataRecord

	for _, record := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
		if validateEmail(cleanEmail) && !emailMap[cleanEmail] {
			emailMap[cleanEmail] = true
			record.Email = cleanEmail
			record.Valid = true
			cleanRecords = append(cleanRecords, record)
		}
	}
	return cleanRecords
}

func main() {
	emails := []string{"test@example.com", "TEST@example.com", "invalid", "another@test.org"}
	uniqueEmails := deduplicateEmails(emails)
	fmt.Println("Deduplicated emails:", uniqueEmails)

	records := []DataRecord{
		{1, "user@domain.com", false},
		{2, "USER@domain.com", false},
		{3, "bad-email", false},
		{4, "new@test.net", false},
	}
	cleanRecords := processRecords(records)
	fmt.Printf("Cleaned records: %+v\n", cleanRecords)
}