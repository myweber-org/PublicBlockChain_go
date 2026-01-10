
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
    var unique []DataRecord

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
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func ValidatePhone(phone string) bool {
    return len(phone) >= 10 && len(phone) <= 15
}

func CleanData(records []DataRecord) []DataRecord {
    var cleaned []DataRecord
    for _, record := range records {
        if ValidateEmail(record.Email) && ValidatePhone(record.Phone) {
            cleaned = append(cleaned, record)
        }
    }
    return DeduplicateRecords(cleaned)
}

func main() {
    sampleData := []DataRecord{
        {1, "test@example.com", "1234567890"},
        {2, "invalid-email", "9876543210"},
        {3, "test@example.com", "1234567890"},
        {4, "another@test.org", "5551234567"},
    }

    cleaned := CleanData(sampleData)
    fmt.Printf("Original: %d records\n", len(sampleData))
    fmt.Printf("Cleaned: %d records\n", len(cleaned))
    for _, r := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", r.ID, r.Email, r.Phone)
    }
}