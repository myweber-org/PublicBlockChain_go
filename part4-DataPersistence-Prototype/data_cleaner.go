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
    Age   string
}

func cleanString(s string) string {
    return strings.TrimSpace(strings.ToLower(s))
}

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processCSVFile(inputPath string, outputPath string) error {
    inFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inFile.Close()

    outFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outFile.Close()

    reader := csv.NewReader(inFile)
    writer := csv.NewWriter(outFile)
    defer writer.Flush()

    header, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read header: %w", err)
    }

    if err := writer.Write(header); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read row: %w", err)
        }

        if len(row) < 4 {
            continue
        }

        record := DataRecord{
            ID:    cleanString(row[0]),
            Name:  cleanString(row[1]),
            Email: cleanString(row[2]),
            Age:   cleanString(row[3]),
        }

        if record.ID == "" || record.Name == "" || !validateEmail(record.Email) {
            continue
        }

        cleanedRow := []string{record.ID, record.Name, record.Email, record.Age}
        if err := writer.Write(cleanedRow); err != nil {
            return fmt.Errorf("failed to write row: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := processCSVFile(inputFile, outputFile); err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
}