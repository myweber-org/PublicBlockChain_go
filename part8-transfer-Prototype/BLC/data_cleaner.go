package main

import (
	"strings"
)

func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func TrimSpaces(slice []string) []string {
	result := make([]string, len(slice))
	for i, item := range slice {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func CleanStringSlice(data []string) []string {
	trimmed := TrimSpaces(data)
	cleaned := RemoveDuplicates(trimmed)
	return cleaned
}package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type Record struct {
    ID      int
    Name    string
    Email   string
    Active  bool
    Score   float64
}

func cleanCSV(inputPath, outputPath string) error {
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

    if err := writer.Write(append(header, "Valid")); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    lineNum := 1
    for {
        lineNum++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("line %d: read error: %v\n", lineNum, err)
            continue
        }

        record, validationErr := validateRecord(row)
        isValid := validationErr == nil

        outputRow := append(row, strconv.FormatBool(isValid))
        if err := writer.Write(outputRow); err != nil {
            fmt.Printf("line %d: write error: %v\n", lineNum, err)
            continue
        }

        if !isValid {
            fmt.Printf("line %d: invalid record: %v\n", lineNum, validationErr)
        }
    }

    return nil
}

func validateRecord(fields []string) (Record, error) {
    if len(fields) != 5 {
        return Record{}, fmt.Errorf("expected 5 fields, got %d", len(fields))
    }

    id, err := strconv.Atoi(fields[0])
    if err != nil {
        return Record{}, fmt.Errorf("invalid ID: %w", err)
    }

    name := strings.TrimSpace(fields[1])
    if name == "" {
        return Record{}, fmt.Errorf("name cannot be empty")
    }

    email := strings.TrimSpace(fields[2])
    if !strings.Contains(email, "@") {
        return Record{}, fmt.Errorf("invalid email format")
    }

    active, err := strconv.ParseBool(fields[3])
    if err != nil {
        return Record{}, fmt.Errorf("invalid active flag: %w", err)
    }

    score, err := strconv.ParseFloat(fields[4], 64)
    if err != nil {
        return Record{}, fmt.Errorf("invalid score: %w", err)
    }

    if score < 0 || score > 100 {
        return Record{}, fmt.Errorf("score must be between 0 and 100")
    }

    return Record{
        ID:     id,
        Name:   name,
        Email:  email,
        Active: active,
        Score:  score,
    }, nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := cleanCSV(inputFile, outputFile); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
}