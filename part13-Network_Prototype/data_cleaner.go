
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
    data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
    cleaned := removeDuplicates(data)
    fmt.Printf("Original: %v\n", data)
    fmt.Printf("Cleaned: %v\n", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

func CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	words := strings.Fields(trimmed)
	seen := make(map[string]bool)
	var result []string

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		if !seen[lowerWord] {
			seen[lowerWord] = true
			result = append(result, word)
		}
	}
	return strings.Join(result, " ")
}

func main() {
	testData := "  Apple banana apple   Cherry BANANA   "
	cleaned := CleanString(testData)
	fmt.Printf("Original: '%s'\n", testData)
	fmt.Printf("Cleaned:  '%s'\n", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Name  string
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) {
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range dc.records {
		key := fmt.Sprintf("%d|%s", record.ID, strings.ToLower(record.Email))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}

	dc.records = unique
	return unique
}

func (dc *DataCleaner) ValidateEmails() []DataRecord {
	var valid []DataRecord

	for _, record := range dc.records {
		if strings.Contains(record.Email, "@") && len(record.Name) > 0 {
			valid = append(valid, record)
		}
	}

	dc.records = valid
	return valid
}

func (dc *DataCleaner) GetRecordCount() int {
	return len(dc.records)
}

func (dc *DataCleaner) PrintRecords() {
	for _, record := range dc.records {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}

func main() {
	cleaner := NewDataCleaner()

	cleaner.AddRecord(DataRecord{ID: 1, Email: "user1@example.com", Name: "Alice"})
	cleaner.AddRecord(DataRecord{ID: 2, Email: "user2@example.com", Name: "Bob"})
	cleaner.AddRecord(DataRecord{ID: 1, Email: "user1@example.com", Name: "Alice"})
	cleaner.AddRecord(DataRecord{ID: 3, Email: "invalid-email", Name: "Charlie"})
	cleaner.AddRecord(DataRecord{ID: 4, Email: "user4@example.com", Name: ""})

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())
	cleaner.PrintRecords()

	cleaner.RemoveDuplicates()
	fmt.Printf("\nAfter deduplication: %d\n", cleaner.GetRecordCount())
	cleaner.PrintRecords()

	cleaner.ValidateEmails()
	fmt.Printf("\nAfter validation: %d\n", cleaner.GetRecordCount())
	cleaner.PrintRecords()
}package main

import "fmt"

func RemoveDuplicates(input []string) []string {
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

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}