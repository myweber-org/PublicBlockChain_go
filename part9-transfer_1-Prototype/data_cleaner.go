
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return len(email) > 5
}

func CleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	uniqueRecords := DeduplicateRecords(records)

	for _, record := range uniqueRecords {
		record.Valid = ValidateEmail(record.Email)
		if record.Valid {
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@test.org", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob", "invalid", false},
	}

	cleaned := CleanData(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d valid records\n", len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", r.ID, r.Name, r.Email)
	}
}
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

type Record struct {
	ID   int
	Name string
	Email string
}

func deduplicateRecords(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record
	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		hash := md5.Sum([]byte(key))
		hashStr := hex.EncodeToString(hash[:])
		if !seen[hashStr] {
			seen[hashStr] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return true
}

func cleanData(records []Record) []Record {
	var valid []Record
	for _, record := range records {
		if validateEmail(record.Email) {
			valid = append(valid, record)
		}
	}
	return deduplicateRecords(valid)
}

func main() {
	records := []Record{
		{1, "Alice", "alice@example.com"},
		{2, "Bob", "bob@example.com"},
		{3, "Alice", "alice@example.com"},
		{4, "Charlie", "invalid-email"},
		{5, "David", "david@example.com"},
		{6, "Bob", "bob@example.com"},
	}

	cleaned := cleanData(records)
	fmt.Printf("Original count: %d\n", len(records))
	fmt.Printf("Cleaned count: %d\n", len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", r.ID, r.Name, r.Email)
	}
}