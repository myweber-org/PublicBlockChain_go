
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

func RemoveDuplicates(records []DataRecord) []DataRecord {
    seen := make(map[int]bool)
    var unique []DataRecord
    for _, record := range records {
        if !seen[record.ID] {
            seen[record.ID] = true
            unique = append(unique, record)
        }
    }
    return unique
}

func ValidateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FormatPhoneNumber(phone string) string {
    cleaned := strings.ReplaceAll(phone, " ", "")
    cleaned = strings.ReplaceAll(cleaned, "-", "")
    if len(cleaned) == 10 {
        return fmt.Sprintf("(%s) %s-%s", cleaned[0:3], cleaned[3:6], cleaned[6:10])
    }
    return phone
}

func ProcessRecords(records []DataRecord) []DataRecord {
    var validRecords []DataRecord
    uniqueRecords := RemoveDuplicates(records)
    
    for _, record := range uniqueRecords {
        if ValidateEmail(record.Email) {
            record.Phone = FormatPhoneNumber(record.Phone)
            validRecords = append(validRecords, record)
        }
    }
    return validRecords
}

func main() {
    sampleData := []DataRecord{
        {ID: 1, Email: "test@example.com", Phone: "1234567890"},
        {ID: 2, Email: "invalid-email", Phone: "987-654-3210"},
        {ID: 1, Email: "test@example.com", Phone: "1234567890"},
        {ID: 3, Email: "user@domain.org", Phone: "555 123 4567"},
    }
    
    processed := ProcessRecords(sampleData)
    fmt.Printf("Processed %d valid records\n", len(processed))
    for _, record := range processed {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", record.ID, record.Email, record.Phone)
    }
}