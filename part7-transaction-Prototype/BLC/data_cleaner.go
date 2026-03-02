
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
		if !seen[email] && isValidEmail(email) {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func validateRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	for _, record := range records {
		if record.ID > 0 && isValidEmail(record.Email) {
			record.Valid = true
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func processData(records []DataRecord) ([]DataRecord, []string) {
	validated := validateRecords(records)
	
	emails := []string{}
	for _, record := range validated {
		emails = append(emails, record.Email)
	}
	
	uniqueEmails := deduplicateEmails(emails)
	return validated, uniqueEmails
}

func main() {
	records := []DataRecord{
		{ID: 1, Email: "user@example.com"},
		{ID: 2, Email: "admin@test.org"},
		{ID: 3, Email: "user@example.com"},
		{ID: 0, Email: "invalid"},
		{ID: 4, Email: "test@domain.com"},
	}

	validRecords, uniqueEmails := processData(records)
	
	fmt.Printf("Valid records: %d\n", len(validRecords))
	fmt.Printf("Unique emails: %v\n", uniqueEmails)
}