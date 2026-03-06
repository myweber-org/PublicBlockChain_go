
package main

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
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
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

func deduplicateEmails(records []DataRecord) []DataRecord {
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

func validateEmails(records []DataRecord) []DataRecord {
	for i := range records {
		records[i].Valid = strings.Contains(records[i].Email, "@") &&
			len(records[i].Email) > 3
	}
	return records
}

func processRecords(records []DataRecord) []DataRecord {
	records = deduplicateEmails(records)
	return validateEmails(records)
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
	}

	processed := processRecords(sampleData)

	for _, record := range processed {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n",
			record.ID, record.Email, record.Valid)
	}
}package main

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
	input := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(input)
	fmt.Printf("Original: %v\n", input)
	fmt.Printf("Cleaned: %v\n", cleaned)
}
package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	// Remove extra whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized := spaceRegex.ReplaceAllString(trimmed, " ")
	
	// Remove non-printable characters
	var result strings.Builder
	for _, r := range normalized {
		if unicode.IsPrint(r) {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

func NormalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

func TruncateString(input string, maxLength int) string {
	if len(input) <= maxLength {
		return input
	}
	
	if maxLength <= 3 {
		return input[:maxLength]
	}
	
	return input[:maxLength-3] + "..."
}