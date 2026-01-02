
package main

import (
    "fmt"
    "strings"
)

// DataCleaner provides methods for cleaning data sets
type DataCleaner struct{}

// RemoveDuplicates removes duplicate strings from a slice
func (dc DataCleaner) RemoveDuplicates(input []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    
    for _, item := range input {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    return result
}

// ValidateEmail checks if a string is a valid email format
func (dc DataCleaner) ValidateEmail(email string) bool {
    if len(email) < 3 || !strings.Contains(email, "@") {
        return false
    }
    
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    
    return len(parts[0]) > 0 && len(parts[1]) > 0 && strings.Contains(parts[1], ".")
}

// TrimSpaces removes leading and trailing whitespace from all strings
func (dc DataCleaner) TrimSpaces(input []string) []string {
    result := make([]string, len(input))
    for i, item := range input {
        result[i] = strings.TrimSpace(item)
    }
    return result
}

func main() {
    cleaner := DataCleaner{}
    
    // Example usage
    data := []string{"  test@example.com", "duplicate@test.com", "test@example.com  ", "invalid", "duplicate@test.com"}
    
    fmt.Println("Original data:", data)
    
    trimmed := cleaner.TrimSpaces(data)
    fmt.Println("After trimming:", trimmed)
    
    deduplicated := cleaner.RemoveDuplicates(trimmed)
    fmt.Println("After deduplication:", deduplicated)
    
    fmt.Println("\nEmail validation results:")
    for _, email := range deduplicated {
        isValid := cleaner.ValidateEmail(email)
        fmt.Printf("%s: %v\n", email, isValid)
    }
}