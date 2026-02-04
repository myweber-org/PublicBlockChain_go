
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
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] && email != "" {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	for i := range records {
		email := records[i].Email
		records[i].Valid = strings.Contains(email, "@") && 
			strings.Contains(email, ".") && 
			len(email) > 5
	}
	return records
}

func CleanDataPipeline(records []DataRecord) []DataRecord {
	records = RemoveDuplicates(records)
	records = ValidateEmails(records)
	return records
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
		{5, "ANOTHER@TEST.ORG", false},
	}

	cleaned := CleanDataPipeline(sampleData)
	
	for _, record := range cleaned {
		status := "Invalid"
		if record.Valid {
			status = "Valid"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", 
			record.ID, record.Email, status)
	}
}
package main

import (
	"encoding/csv"
	"fmt"
	"strings"
)

func TrimCSVColumns(records [][]string) [][]string {
	trimmed := make([][]string, len(records))
	for i, row := range records {
		trimmed[i] = make([]string, len(row))
		for j, val := range row {
			trimmed[i][j] = strings.TrimSpace(val)
		}
	}
	return trimmed
}

func main() {
	data := [][]string{
		{"  id  ", " name ", " value "},
		{"  1", "alpha  ", "  100"},
		{"2  ", "  beta", "200  "},
	}

	cleaned := TrimCSVColumns(data)
	writer := csv.NewWriter(fmt.Stdout)
	writer.WriteAll(cleaned)
	writer.Flush()
}