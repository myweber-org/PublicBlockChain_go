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
		if !seen[email] && email != "" {
			seen[email] = true
			unique = append(unique, DataRecord{
				ID:    record.ID,
				Email: email,
				Valid: validateEmail(email),
			})
		}
	}
	return unique
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func printRecords(records []DataRecord) {
	for _, r := range records {
		status := "INVALID"
		if r.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", r.ID, r.Email, status)
	}
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "test@domain.org", false},
		{4, "invalid-email", false},
		{5, "user@example.com", false},
		{6, "", false},
	}

	cleaned := deduplicateEmails(sampleData)
	fmt.Println("Cleaned Data Records:")
	printRecords(cleaned)
}