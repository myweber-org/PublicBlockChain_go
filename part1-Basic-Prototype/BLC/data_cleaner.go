
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
}package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput cleans and normalizes user-provided strings
func SanitizeInput(input string) string {
	// Remove leading/trailing whitespace
	cleaned := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned = spaceRegex.ReplaceAllString(cleaned, " ")
	
	// Remove potentially dangerous characters (customize as needed)
	dangerousChars := regexp.MustCompile(`[<>{}]`)
	cleaned = dangerousChars.ReplaceAllString(cleaned, "")
	
	return cleaned
}

// NormalizeEmail formats email addresses to lowercase and trims spaces
func NormalizeEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	
	// Basic email validation pattern
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ""
	}
	
	return email
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

func DeduplicateEmails(emails []string) []string {
    seen := make(map[string]struct{})
    result := []string{}
    for _, email := range emails {
        if _, exists := seen[email]; !exists {
            seen[email] = struct{}{}
            result = append(result, email)
        }
    }
    return result
}

func ValidateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
    emailSet := make(map[string]bool)
    cleaned := []DataRecord{}
    
    for _, record := range records {
        if ValidateEmail(record.Email) && !emailSet[record.Email] {
            emailSet[record.Email] = true
            record.Valid = true
            cleaned = append(cleaned, record)
        }
    }
    return cleaned
}

func main() {
    sampleData := []DataRecord{
        {1, "user@example.com", false},
        {2, "invalid-email", false},
        {3, "user@example.com", false},
        {4, "test@domain.org", false},
    }
    
    cleaned := CleanData(sampleData)
    fmt.Printf("Cleaned records: %d\n", len(cleaned))
    for _, r := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
    }
    
    emails := []string{"a@b.com", "c@d.com", "a@b.com", "e@f.com"}
    unique := DeduplicateEmails(emails)
    fmt.Println("Unique emails:", unique)
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
	numbers := []int{1, 2, 2, 3, 4, 4, 5, 6, 6, 7}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)
}
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
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange", "banana"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
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

func DeduplicateEmails(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			record.Email = email
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmailFormat(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
	records = DeduplicateEmails(records)
	for i := range records {
		records[i].Valid = ValidateEmailFormat(records[i].Email)
	}
	return records
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
	}

	cleaned := CleanData(sampleData)
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", record.ID, record.Email, record.Valid)
	}
}