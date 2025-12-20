
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