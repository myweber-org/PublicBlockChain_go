
package main

import (
    "fmt"
    "strings"
)

type DataCleaner struct {
    duplicates map[string]bool
}

func NewDataCleaner() *DataCleaner {
    return &DataCleaner{
        duplicates: make(map[string]bool),
    }
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
    unique := []string{}
    for _, item := range items {
        normalized := strings.ToLower(strings.TrimSpace(item))
        if !dc.duplicates[normalized] {
            dc.duplicates[normalized] = true
            unique = append(unique, item)
        }
    }
    return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    return len(parts[0]) > 0 && len(parts[1]) > 0
}

func main() {
    cleaner := NewDataCleaner()
    
    data := []string{"apple", "banana", "Apple", "cherry", "banana", "  BANANA  "}
    unique := cleaner.RemoveDuplicates(data)
    fmt.Println("Unique items:", unique)
    
    emails := []string{"test@example.com", "invalid-email", "user@domain", "@nodomain", "domain@", ""}
    for _, email := range emails {
        fmt.Printf("Email '%s' valid: %v\n", email, cleaner.ValidateEmail(email))
    }
}