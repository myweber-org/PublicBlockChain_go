
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataProcessor struct {
    inputPath  string
    outputPath string
    delimiter  rune
}

func NewDataProcessor(input, output string) *DataProcessor {
    return &DataProcessor{
        inputPath:  input,
        outputPath: output,
        delimiter:  ',',
    }
}

func (dp *DataProcessor) SetDelimiter(delim rune) {
    dp.delimiter = delim
}

func (dp *DataProcessor) Process() error {
    inputFile, err := os.Open(dp.inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(dp.outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    reader := csv.NewReader(inputFile)
    reader.Comma = dp.delimiter
    writer := csv.NewWriter(outputFile)
    writer.Comma = dp.delimiter
    defer writer.Flush()

    lineCount := 0
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV: %w", err)
        }

        cleanedRecord := dp.cleanRecord(record)
        if len(cleanedRecord) > 0 {
            if err := writer.Write(cleanedRecord); err != nil {
                return fmt.Errorf("error writing CSV: %w", err)
            }
            lineCount++
        }
    }

    fmt.Printf("Processed %d lines successfully\n", lineCount)
    return nil
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
    cleaned := make([]string, 0, len(record))
    for _, field := range record {
        cleanedField := strings.TrimSpace(field)
        if cleanedField != "" {
            cleaned = append(cleaned, cleanedField)
        }
    }
    return cleaned
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: data_processor <input.csv> <output.csv>")
        os.Exit(1)
    }

    processor := NewDataProcessor(os.Args[1], os.Args[2])
    if err := processor.Process(); err != nil {
        fmt.Printf("Processing error: %v\n", err)
        os.Exit(1)
    }
}