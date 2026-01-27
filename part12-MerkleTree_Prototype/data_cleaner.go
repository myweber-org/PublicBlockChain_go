
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
	if len(record) == 0 {
		return false
	}
	if strings.Contains(record, "test") {
		return false
	}
	return true
}

func (dc *DataCleaner) Reset() {
	dc.processedRecords = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	sampleData := []string{
		"Customer A",
		"customer a",
		"Customer B",
		"test record",
		"",
		"Customer C",
		"  Customer A  ",
	}
	
	fmt.Println("Original records:", sampleData)
	cleaned := cleaner.RemoveDuplicates(sampleData)
	fmt.Println("Cleaned records:", cleaned)
	
	cleaner.Reset()
	fmt.Println("Cleaner reset completed")
}
package main

import (
	"fmt"
)

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