
package data

import (
	"errors"
	"regexp"
	"strings"
)

type Record struct {
	ID    string
	Email string
	Score int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateRecord(r Record) error {
	if r.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}
	if r.Score < 0 || r.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func TransformRecords(records []Record) ([]Record, error) {
	var processed []Record
	for _, r := range records {
		if err := ValidateRecord(r); err != nil {
			return nil, err
		}
		r.Email = NormalizeEmail(r.Email)
		processed = append(processed, r)
	}
	return processed, nil
}

func CalculateAverage(records []Record) float64 {
	if len(records) == 0 {
		return 0.0
	}
	total := 0
	for _, r := range records {
		total += r.Score
	}
	return float64(total) / float64(len(records))
}
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

func processCSVFile(inputPath string, outputPath string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    reader := csv.NewReader(inputFile)
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    lineNumber := 0
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV at line %d: %w", lineNumber, err)
        }

        lineNumber++
        if lineNumber == 1 {
            if err := writer.Write(record); err != nil {
                return fmt.Errorf("error writing header: %w", err)
            }
            continue
        }

        cleanedRecord := cleanRecord(record)
        if cleanedRecord == nil {
            continue
        }

        if err := writer.Write(cleanedRecord); err != nil {
            return fmt.Errorf("error writing record at line %d: %w", lineNumber, err)
        }
    }

    return nil
}

func cleanRecord(record []string) []string {
    cleaned := make([]string, len(record))
    for i, field := range record {
        cleaned[i] = strings.TrimSpace(field)
        if cleaned[i] == "" {
            return nil
        }
    }
    return cleaned
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_processor <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := processCSVFile(inputFile, outputFile); err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %s to %s\n", inputFile, outputFile)
}