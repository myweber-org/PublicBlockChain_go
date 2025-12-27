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

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	for _, record := range records {
		if !seen[record.Email] {
			seen[record.Email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	for i := range records {
		records[i].Valid = strings.Contains(records[i].Email, "@") && 
			strings.Contains(records[i].Email, ".")
	}
	return records
}

func CleanData(records []DataRecord) []DataRecord {
	unique := RemoveDuplicates(records)
	validated := ValidateEmails(unique)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "invalid-email", false},
		{3, "user@example.com", false},
		{4, "test@domain.org", false},
	}

	cleaned := CleanData(sampleData)
	
	for _, record := range cleaned {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", 
			record.ID, record.Email, status)
	}
}