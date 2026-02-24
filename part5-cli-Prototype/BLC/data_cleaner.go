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

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
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
		if !dc.processedRecords[normalized] && dc.isValidRecord(normalized) {
			dc.processedRecords[normalized] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func (dc *DataCleaner) isValidRecord(record string) bool {
	return len(record) > 0 && !strings.ContainsAny(record, "!@#$%")
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{
		"user@example.com",
		"user@example.com",
		"test@domain.org",
		"invalid-email",
		"another@test.com",
		"",
	}
	
	fmt.Println("Original records:", len(data))
	uniqueData := cleaner.RemoveDuplicates(data)
	fmt.Println("Unique records:", len(uniqueData))
	
	for _, record := range uniqueData {
		if cleaner.ValidateEmail(record) {
			fmt.Printf("Valid email: %s\n", record)
		} else {
			fmt.Printf("Invalid record: %s\n", record)
		}
	}
}