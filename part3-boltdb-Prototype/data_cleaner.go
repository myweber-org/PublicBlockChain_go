
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

    cleaner.AddRecord(DataRecord{ID: 1, Email: "user@example.com", Phone: "1234567890"})
    cleaner.AddRecord(DataRecord{ID: 2, Email: "user@example.com", Phone: "1234567890"})
    cleaner.AddRecord(DataRecord{ID: 3, Email: "invalid-email", Phone: "0987654321"})
    cleaner.AddRecord(DataRecord{ID: 4, Email: "another@test.org", Phone: "5551234567"})

    fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

    unique := cleaner.RemoveDuplicates()
    fmt.Printf("After deduplication: %d\n", len(unique))

    valid := cleaner.ValidateEmails()
    fmt.Printf("Valid email records: %d\n", len(valid))

    for _, record := range valid {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", record.ID, record.Email, record.Phone)
    }
}package main

import (
	"regexp"
	"strings"
)

func SanitizeCSVField(input string) string {
	if input == "" {
		return input
	}

	// Remove leading/trailing whitespace
	trimmed := strings.TrimSpace(input)

	// Escape double quotes by doubling them (standard CSV escaping)
	escaped := strings.ReplaceAll(trimmed, `"`, `""`)

	// Check if field needs quoting
	needsQuoting := false
	if strings.ContainsAny(escaped, `,"`) {
		needsQuoting = true
	} else {
		// Check for newlines or carriage returns
		re := regexp.MustCompile(`[\r\n]`)
		if re.MatchString(escaped) {
			needsQuoting = true
		}
	}

	if needsQuoting {
		return `"` + escaped + `"`
	}
	return escaped
}package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
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