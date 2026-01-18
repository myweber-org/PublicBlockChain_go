
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
    unique := []DataRecord{}

    for _, record := range records {
        normalizedEmail := strings.ToLower(strings.TrimSpace(record.Email))
        if !seen[normalizedEmail] {
            seen[normalizedEmail] = true
            unique = append(unique, record)
        }
    }
    return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
    validated := []DataRecord{}
    for _, record := range records {
        record.Valid = strings.Contains(record.Email, "@") && len(record.Email) > 3
        validated = append(validated, record)
    }
    return validated
}

func PrintRecords(records []DataRecord) {
    for _, record := range records {
        status := "INVALID"
        if record.Valid {
            status = "VALID"
        }
        fmt.Printf("ID: %d, Email: %s, Status: %s\n", record.ID, record.Email, status)
    }
}

func main() {
    sampleData := []DataRecord{
        {1, "user@example.com", false},
        {2, "user@example.com", false},
        {3, "admin@test.org", false},
        {4, "invalid-email", false},
        {5, "  USER@EXAMPLE.COM  ", false},
    }

    fmt.Println("Original data:")
    PrintRecords(sampleData)

    uniqueData := RemoveDuplicates(sampleData)
    fmt.Println("\nAfter deduplication:")
    PrintRecords(uniqueData)

    validatedData := ValidateEmails(uniqueData)
    fmt.Println("\nAfter validation:")
    PrintRecords(validatedData)
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}