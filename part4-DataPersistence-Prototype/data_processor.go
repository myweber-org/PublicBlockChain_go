
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID      string
    Name    string
    Email   string
    Active  string
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    var records []DataRecord
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if lineNumber == 1 {
            continue
        }

        if len(row) < 4 {
            return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
        }

        record := DataRecord{
            ID:     strings.TrimSpace(row[0]),
            Name:   strings.TrimSpace(row[1]),
            Email:  strings.TrimSpace(row[2]),
            Active: strings.TrimSpace(row[3]),
        }

        if record.ID == "" || record.Name == "" {
            return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
        }

        records = append(records, record)
    }

    return records, nil
}

func ValidateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FilterActiveUsers(records []DataRecord) []DataRecord {
    var activeUsers []DataRecord
    for _, record := range records {
        if record.Active == "true" && ValidateEmail(record.Email) {
            activeUsers = append(activeUsers, record)
        }
    }
    return activeUsers
}

func GenerateReport(records []DataRecord) {
    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Println("Active users with valid emails:")
    for _, user := range FilterActiveUsers(records) {
        fmt.Printf("  - %s (%s)\n", user.Name, user.Email)
    }
}