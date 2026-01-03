
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type Record struct {
    ID        int
    Name      string
    Email     string
    Age       int
    Validated bool
}

func cleanEmail(email string) string {
    return strings.ToLower(strings.TrimSpace(email))
}

func validateAge(age int) bool {
    return age >= 18 && age <= 120
}

func processCSVFile(inputPath string) ([]Record, error) {
    file, err := os.Open(inputPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
    lineNumber := 0

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        lineNumber++
        if lineNumber == 1 {
            continue
        }

        if len(row) < 4 {
            continue
        }

        id, err := strconv.Atoi(strings.TrimSpace(row[0]))
        if err != nil {
            continue
        }

        name := strings.TrimSpace(row[1])
        email := cleanEmail(row[2])
        age, err := strconv.Atoi(strings.TrimSpace(row[3]))
        if err != nil {
            continue
        }

        record := Record{
            ID:        id,
            Name:      name,
            Email:     email,
            Age:       age,
            Validated: validateAge(age),
        }

        records = append(records, record)
    }

    return records, nil
}

func generateReport(records []Record) {
    validCount := 0
    for _, record := range records {
        if record.Validated {
            validCount++
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Valid records: %d\n", validCount)
    fmt.Printf("Invalid records: %d\n", len(records)-validCount)
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run data_cleaner.go <input_file.csv>")
        return
    }

    records, err := processCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        return
    }

    generateReport(records)
}