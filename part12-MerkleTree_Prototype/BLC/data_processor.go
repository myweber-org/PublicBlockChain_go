
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
}

func NewDataProcessor(input, output string) *DataProcessor {
    return &DataProcessor{
        inputPath:  input,
        outputPath: output,
    }
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
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read CSV record: %w", err)
        }

        cleanedRecord := dp.cleanRecord(record)
        if dp.isValidRecord(cleanedRecord) {
            if err := writer.Write(cleanedRecord); err != nil {
                return fmt.Errorf("failed to write CSV record: %w", err)
            }
        }
    }

    return nil
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
    cleaned := make([]string, len(record))
    for i, field := range record {
        cleaned[i] = strings.TrimSpace(field)
    }
    return cleaned
}

func (dp *DataProcessor) isValidRecord(record []string) bool {
    for _, field := range record {
        if field == "" {
            return false
        }
    }
    return len(record) > 0
}

func main() {
    processor := NewDataProcessor("input.csv", "output.csv")
    if err := processor.Process(); err != nil {
        fmt.Printf("Processing error: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Data processing completed successfully")
}