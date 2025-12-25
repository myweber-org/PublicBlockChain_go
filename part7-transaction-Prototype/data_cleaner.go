
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
	dangerous := []string{"<", ">", "'", "\"", "&"}
	sanitized := input
	for _, char := range dangerous {
		sanitized = strings.ReplaceAll(sanitized, char, "")
	}
	return strings.TrimSpace(sanitized)
}

func main() {
	cleaner := NewDataCleaner()
	
	records := []string{"test@example.com", "TEST@example.com", "invalid", "user@domain"}
	fmt.Println("Original:", records)
	
	deduped := cleaner.RemoveDuplicates(records)
	fmt.Println("Deduplicated:", deduped)
	
	for _, email := range deduped {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("%s is valid\n", email)
		} else {
			fmt.Printf("%s is invalid\n", email)
		}
	}
	
	sample := "<script>alert('test')</script>"
	fmt.Println("Sanitized:", cleaner.SanitizeInput(sample))
}
package main

import "fmt"

func removeDuplicates(input []int) []int {
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
	data := []int{4, 2, 8, 2, 4, 9, 8, 1}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
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
package main

import "fmt"

func removeDuplicates(nums []int) []int {
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
    input := []int{1, 2, 2, 3, 4, 4, 5}
    output := removeDuplicates(input)
    fmt.Println("Original:", input)
    fmt.Println("Cleaned:", output)
}