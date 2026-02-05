package main

import (
    "fmt"
    "strings"
)

type DataRecord struct {
    ID    int
    Email string
    Phone string
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
    unique := make([]DataRecord, 0)

    for _, record := range dc.records {
        key := fmt.Sprintf("%s|%s", record.Email, record.Phone)
        if !seen[key] {
            seen[key] = true
            unique = append(unique, record)
        }
    }

    dc.records = unique
    return unique
}

func (dc *DataCleaner) ValidateEmails() []DataRecord {
    valid := make([]DataRecord, 0)

    for _, record := range dc.records {
        if strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".") {
            valid = append(valid, record)
        }
    }

    return valid
}

func (dc *DataCleaner) GetRecordCount() int {
    return len(dc.records)
}

func main() {
    cleaner := NewDataCleaner()

    sampleData := []DataRecord{
        {1, "user@example.com", "1234567890"},
        {2, "user@example.com", "1234567890"},
        {3, "invalid-email", "0987654321"},
        {4, "another@test.org", "5551234567"},
    }

    for _, record := range sampleData {
        cleaner.AddRecord(record)
    }

    fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())
    
    unique := cleaner.RemoveDuplicates()
    fmt.Printf("After deduplication: %d\n", len(unique))
    
    valid := cleaner.ValidateEmails()
    fmt.Printf("Valid email records: %d\n", len(valid))
}
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

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"Apple", "apple", " BANANA ", "banana", "Cherry"}
	fmt.Println("Original:", data)
	
	deduped := cleaner.Deduplicate(data)
	fmt.Println("Deduplicated:", deduped)
	
	cleaner.Reset()
	testValue := "  TEST  "
	fmt.Printf("Normalized '%s': '%s'\n", testValue, cleaner.Normalize(testValue))
}