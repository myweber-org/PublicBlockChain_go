
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID    string
    Name  string
    Email string
    Valid bool
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

        if len(row) < 3 {
            continue
        }

        record := DataRecord{
            ID:    strings.TrimSpace(row[0]),
            Name:  strings.TrimSpace(row[1]),
            Email: strings.TrimSpace(row[2]),
            Valid: validateEmail(strings.TrimSpace(row[2])),
        }

        if record.ID != "" && record.Name != "" {
            records = append(records, record)
        }
    }

    return records, nil
}

func validateEmail(email string) bool {
    if email == "" {
        return false
    }
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func GenerateReport(records []DataRecord) {
    validCount := 0
    for _, record := range records {
        if record.Valid {
            validCount++
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Valid email addresses: %d\n", validCount)
    fmt.Printf("Invalid email addresses: %d\n", len(records)-validCount)

    if len(records) > 0 {
        fmt.Println("\nFirst 5 records:")
        for i := 0; i < len(records) && i < 5; i++ {
            status := "INVALID"
            if records[i].Valid {
                status = "VALID"
            }
            fmt.Printf("%s: %s <%s> [%s]\n", 
                records[i].ID, 
                records[i].Name, 
                records[i].Email, 
                status)
        }
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run data_processor.go <csv_file>")
        os.Exit(1)
    }

    filename := os.Args[1]
    records, err := ProcessCSVFile(filename)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    GenerateReport(records)
}