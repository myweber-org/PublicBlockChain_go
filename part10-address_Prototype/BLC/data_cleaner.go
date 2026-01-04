
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

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(normalized) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && !strings.ContainsAny(item, "!@#$%")
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	data := []string{"apple", "Apple", "banana", "", "cherry!", "banana", "date "}
	
	cleaned := cleaner.RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
	
	cleaner.Reset()
	testData := []string{"test1", "test2", "test1"}
	fmt.Printf("Reset test: %v\n", cleaner.RemoveDuplicates(testData))
}
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

func DeduplicateRecords(records []DataRecord) []DataRecord {
    seen := make(map[string]bool)
    unique := []DataRecord{}
    
    for _, record := range records {
        key := fmt.Sprintf("%s|%s", record.Email, record.Phone)
        if !seen[key] {
            seen[key] = true
            unique = append(unique, record)
        }
    }
    return unique
}

func ValidateEmail(email string) bool {
    if len(email) == 0 {
        return false
    }
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func ValidatePhone(phone string) bool {
    if len(phone) == 0 {
        return false
    }
    for _, ch := range phone {
        if ch < '0' || ch > '9' {
            return false
        }
    }
    return len(phone) >= 10
}

func CleanData(records []DataRecord) []DataRecord {
    validRecords := []DataRecord{}
    
    for _, record := range records {
        if ValidateEmail(record.Email) && ValidatePhone(record.Phone) {
            validRecords = append(validRecords, record)
        }
    }
    
    return DeduplicateRecords(validRecords)
}

func main() {
    sampleData := []DataRecord{
        {1, "test@example.com", "1234567890"},
        {2, "test@example.com", "1234567890"},
        {3, "invalid-email", "1234567890"},
        {4, "another@test.com", "not-a-phone"},
        {5, "valid@email.org", "0987654321"},
    }
    
    cleaned := CleanData(sampleData)
    
    fmt.Printf("Original records: %d\n", len(sampleData))
    fmt.Printf("Cleaned records: %d\n", len(cleaned))
    
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", 
            record.ID, record.Email, record.Phone)
    }
}