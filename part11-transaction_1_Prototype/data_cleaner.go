package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSVData(inputPath, outputPath string) error {
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

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	trimmedHeaders := make([]string, len(headers))
	for i, h := range headers {
		trimmedHeaders[i] = strings.TrimSpace(h)
	}

	if err := writer.Write(trimmedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		cleanedRecord := make([]string, len(record))
		hasEmpty := false
		for i, field := range record {
			cleaned := strings.TrimSpace(field)
			if cleaned == "" {
				hasEmpty = true
			}
			cleanedRecord[i] = cleaned
		}

		if !hasEmpty {
			if err := writer.Write(cleanedRecord); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
			recordCount++
		}
	}

	fmt.Printf("Cleaned %d records, removed empty rows\n", recordCount)
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSVData(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data cleaning completed successfully")
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
	ID    int
	Name  string
	Email string
	Score float64
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

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	lineNum := 1
	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading row %d: %w", lineNum, err)
		}

		record, err := validateRow(row)
		if err != nil {
			fmt.Printf("Skipping row %d: %v\n", lineNum, err)
			continue
		}

		cleanedRow := []string{
			strconv.Itoa(record.ID),
			strings.TrimSpace(record.Name),
			strings.ToLower(strings.TrimSpace(record.Email)),
			fmt.Sprintf("%.2f", record.Score),
		}

		if err := writer.Write(cleanedRow); err != nil {
			return fmt.Errorf("failed to write row %d: %w", lineNum, err)
		}
	}

	return nil
}

func validateRow(row []string) (*Record, error) {
	if len(row) != 4 {
		return nil, fmt.Errorf("expected 4 columns, got %d", len(row))
	}

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	if id <= 0 {
		return nil, fmt.Errorf("ID must be positive")
	}

	name := strings.TrimSpace(row[1])
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	email := strings.TrimSpace(row[2])
	if !strings.Contains(email, "@") {
		return nil, fmt.Errorf("invalid email format")
	}

	score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid score: %w", err)
	}
	if score < 0 || score > 100 {
		return nil, fmt.Errorf("score must be between 0 and 100")
	}

	return &Record{
		ID:    id,
		Name:  name,
		Email: email,
		Score: score,
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

	fmt.Printf("Successfully cleaned data. Output written to %s\n", outputFile)
}