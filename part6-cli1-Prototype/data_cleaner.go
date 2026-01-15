
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    
    for _, num := range nums {
        if !seen[num] {
            seen[num] = true
            result = append(result, num)
        }
    }
    
    return result
}

func main() {
    data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
    cleaned := RemoveDuplicates(data)
    fmt.Printf("Original: %v\n", data)
    fmt.Printf("Cleaned: %v\n", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func CleanData(data []string) []string {
	cleaned := []string{}
	for _, item := range data {
		normalized := NormalizeString(item)
		cleaned = append(cleaned, normalized)
	}
	return RemoveDuplicates(cleaned)
}

func main() {
	rawData := []string{"  Apple", "banana  ", "APPLE", "Banana", "  cherry  "}
	cleanedData := CleanData(rawData)
	fmt.Println("Cleaned data:", cleanedData)
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

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if strings.Contains(email, "@") && strings.Contains(email, ".") {
			record.Valid = true
			valid = append(valid, record)
		}
	}
	return valid
}

func CleanData(records []DataRecord) []DataRecord {
	unique := RemoveDuplicates(records)
	valid := ValidateEmails(unique)
	return valid
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "test@domain.org", false},
		{4, "invalid-email", false},
		{5, "another@test.co", false},
	}

	cleaned := CleanData(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
	
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", record.ID, record.Email, record.Valid)
	}
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
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
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
	validRecords := []DataRecord{}
	for _, record := range records {
		if ValidateEmail(record.Email) {
			record.Valid = true
			validRecords = append(validRecords, record)
		}
	}
	return DeduplicateRecords(validRecords)
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "invalid-email", false},
		{3, "user@example.com", false},
		{4, "test@domain.org", false},
	}

	cleaned := CleanData(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}