
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	Data []string
}

func NewDataCleaner(data []string) *DataCleaner {
	return &DataCleaner{Data: data}
}

func (dc *DataCleaner) RemoveDuplicates() []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range dc.Data {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; !exists {
			seen[trimmed] = struct{}{}
			result = append(result, trimmed)
		}
	}
	return result
}

func (dc *DataCleaner) Clean() []string {
	return dc.RemoveDuplicates()
}

func main() {
	rawData := []string{"  apple ", "banana", "  apple", "cherry  ", "", "banana", "date"}
	cleaner := NewDataCleaner(rawData)
	cleaned := cleaner.Clean()
	fmt.Println("Cleaned data:", cleaned)
}
package main

import (
    "fmt"
    "strings"
)

type DataRecord struct {
    ID    int
    Name  string
    Email string
    Valid bool
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
    seen := make(map[string]bool)
    unique := []DataRecord{}
    
    for _, record := range records {
        key := fmt.Sprintf("%s|%s", record.Name, record.Email)
        if !seen[key] {
            seen[key] = true
            unique = append(unique, record)
        }
    }
    return unique
}

func ValidateEmail(email string) bool {
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    return len(parts[0]) > 0 && len(parts[1]) > 0
}

func ValidateRecords(records []DataRecord) []DataRecord {
    validated := []DataRecord{}
    for _, record := range records {
        record.Valid = ValidateEmail(record.Email)
        validated = append(validated, record)
    }
    return validated
}

func CleanData(records []DataRecord) []DataRecord {
    unique := DeduplicateRecords(records)
    return ValidateRecords(unique)
}

func main() {
    sampleData := []DataRecord{
        {1, "John Doe", "john@example.com", false},
        {2, "Jane Smith", "jane@example.com", false},
        {3, "John Doe", "john@example.com", false},
        {4, "Bob Wilson", "invalid-email", false},
    }
    
    cleaned := CleanData(sampleData)
    
    fmt.Printf("Original records: %d\n", len(sampleData))
    fmt.Printf("Cleaned records: %d\n", len(cleaned))
    
    for _, record := range cleaned {
        status := "INVALID"
        if record.Valid {
            status = "VALID"
        }
        fmt.Printf("ID: %d, Name: %s, Status: %s\n", 
            record.ID, record.Name, status)
    }
}