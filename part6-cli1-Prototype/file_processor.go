
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID    int
    Name  string
    Value float64
}

func processCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
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

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
        }

        name := row[1]
        if name == "" {
            return nil, fmt.Errorf("empty name at line %d", lineNumber)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, fmt.Errorf("no valid records found in file")
    }

    return records, nil
}

func validateRecords(records []Record) error {
    seenIDs := make(map[int]bool)
    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("record ID must be positive: %d", record.ID)
        }
        if seenIDs[record.ID] {
            return fmt.Errorf("duplicate ID found: %d", record.ID)
        }
        seenIDs[record.ID] = true

        if record.Value < 0 {
            return fmt.Errorf("negative value not allowed for record %d", record.ID)
        }
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: file_processor <csv_file>")
        os.Exit(1)
    }

    filename := os.Args[1]
    records, err := processCSVFile(filename)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    if err := validateRecords(records); err != nil {
        fmt.Printf("Validation error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %d records\n", len(records))
    totalValue := 0.0
    for _, record := range records {
        totalValue += record.Value
        fmt.Printf("ID: %d, Name: %s, Value: %.2f\n", record.ID, record.Name, record.Value)
    }
    fmt.Printf("Total value: %.2f\n", totalValue)
}