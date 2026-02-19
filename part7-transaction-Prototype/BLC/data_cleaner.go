
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) Deduplicate(values []string) []string {
	dc.seen = make(map[string]bool)
	var result []string
	for _, v := range values {
		if !dc.IsDuplicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"Apple", "apple ", " BANANA", "banana", "Cherry"}
	fmt.Println("Original:", data)
	
	deduped := cleaner.Deduplicate(data)
	fmt.Println("Deduplicated:", deduped)
	
	testValue := "  APPLE  "
	fmt.Printf("Is '%s' duplicate? %v\n", testValue, cleaner.IsDuplicate(testValue))
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
	numbers := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Printf("Original: %v\n", numbers)
	fmt.Printf("Cleaned: %v\n", uniqueNumbers)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedCount int
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{processedCount: 0}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" && !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	dc.processedCount += len(items)
	return result
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0 && strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) GetStats() string {
	return fmt.Sprintf("Processed %d items", dc.processedCount)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"  john@example.com", "john@example.com", "invalid-email", "  ", "sarah@test.org"}
	
	uniqueEmails := cleaner.RemoveDuplicates(data)
	fmt.Println("Unique emails:", uniqueEmails)
	
	for _, email := range uniqueEmails {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("%s is valid\n", email)
		} else {
			fmt.Printf("%s is invalid\n", email)
		}
	}
	
	fmt.Println(cleaner.GetStats())
}