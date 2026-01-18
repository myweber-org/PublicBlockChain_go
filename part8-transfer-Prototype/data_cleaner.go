
package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange", "banana"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}
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
    result := []DataRecord{}
    for _, record := range records {
        normalizedEmail := strings.ToLower(strings.TrimSpace(record.Email))
        if !seen[normalizedEmail] {
            seen[normalizedEmail] = true
            result = append(result, record)
        }
    }
    return result
}

func ValidateEmails(records []DataRecord) []DataRecord {
    for i := range records {
        email := records[i].Email
        records[i].Valid = strings.Contains(email, "@") && strings.Contains(email, ".")
    }
    return records
}

func main() {
    sampleData := []DataRecord{
        {1, "user@example.com", false},
        {2, "USER@example.com", false},
        {3, "test.user@domain.org", false},
        {4, "invalid-email", false},
        {5, "user@example.com", false},
    }

    fmt.Println("Original records:", len(sampleData))
    unique := RemoveDuplicates(sampleData)
    fmt.Println("After deduplication:", len(unique))
    
    validated := ValidateEmails(unique)
    validCount := 0
    for _, r := range validated {
        if r.Valid {
            validCount++
        }
    }
    fmt.Println("Valid email records:", validCount)
}