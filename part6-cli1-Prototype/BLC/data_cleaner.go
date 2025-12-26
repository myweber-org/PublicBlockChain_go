
package main

import (
    "fmt"
    "strings"
)

// DataCleaner provides methods for cleaning datasets
type DataCleaner struct{}

// RemoveDuplicates removes duplicate entries from a slice of strings
func (dc *DataCleaner) RemoveDuplicates(data []string) []string {
    seen := make(map[string]struct{})
    result := []string{}
    for _, item := range data {
        if _, exists := seen[item]; !exists {
            seen[item] = struct{}{}
            result = append(result, item)
        }
    }
    return result
}

// ValidateEmail checks if a string is a valid email format
func (dc *DataCleaner) ValidateEmail(email string) bool {
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
    return strings.Contains(parts[1], ".")
}

// NormalizeWhitespace removes extra spaces from text
func (dc *DataCleaner) NormalizeWhitespace(text string) string {
    words := strings.Fields(text)
    return strings.Join(words, " ")
}

func main() {
    cleaner := &DataCleaner{}
    
    sampleData := []string{"apple", "banana", "apple", "cherry", "banana"}
    uniqueData := cleaner.RemoveDuplicates(sampleData)
    fmt.Println("Deduplicated:", uniqueData)
    
    emails := []string{"test@example.com", "invalid-email", "user@domain"}
    for _, email := range emails {
        fmt.Printf("Email %s valid: %v\n", email, cleaner.ValidateEmail(email))
    }
    
    text := "  This   has   extra   spaces   "
    normalized := cleaner.NormalizeWhitespace(text)
    fmt.Println("Normalized text:", normalized)
}
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
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedRecords map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processedRecords: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(records []string) []string {
	var unique []string
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] {
			dc.processedRecords[normalized] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	replacer := strings.NewReplacer("\n", " ", "\t", " ", "\r", " ")
	return replacer.Replace(trimmed)
}

func main() {
	cleaner := NewDataCleaner()
	
	duplicateData := []string{"user@example.com", "test@domain.org", "user@example.com", "TEST@DOMAIN.ORG"}
	uniqueEmails := cleaner.RemoveDuplicates(duplicateData)
	fmt.Println("Unique emails:", uniqueEmails)
	
	testEmails := []string{"valid@test.com", "invalid", "no@tld", "@missinglocal.com"}
	for _, email := range testEmails {
		fmt.Printf("Email %s valid: %v\n", email, cleaner.ValidateEmail(email))
	}
	
	dirtyInput := "\t  Sample  \n  Text  \r\n  "
	cleanInput := cleaner.SanitizeInput(dirtyInput)
	fmt.Printf("Sanitized: '%s'\n", cleanInput)
}