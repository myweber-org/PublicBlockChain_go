
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
	var unique []string
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
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"test@example.com", "  TEST@example.com  ", "invalid", "another@test.org"}
	fmt.Println("Original:", data)
	
	uniqueData := cleaner.RemoveDuplicates(data)
	fmt.Println("Deduplicated:", uniqueData)
	
	for _, email := range uniqueData {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("Valid email: %s\n", email)
		} else {
			fmt.Printf("Invalid email: %s\n", email)
		}
	}
}