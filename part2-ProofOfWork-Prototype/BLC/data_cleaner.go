package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) Deduplicate(values []string) []string {
	dc.seen = make(map[string]bool)
	var result []string
	for _, v := range values {
		if !dc.IsDuplicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	cleaner := NewDataCleaner()
	data := []string{"apple", "Apple ", "banana", "BANANA", "  cherry  "}
	
	fmt.Println("Original data:", data)
	
	deduped := cleaner.Deduplicate(data)
	fmt.Println("Deduplicated:", deduped)
	
	for _, item := range deduped {
		fmt.Printf("Normalized: '%s'\n", cleaner.Normalize(item))
	}
}
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

	headers = append(headers, "Valid")
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read row: %w", err)
		}

		record, validationErr := validateRow(row)
		valid := "true"
		if validationErr != nil {
			valid = "false"
		}

		outputRow := []string{
			strconv.Itoa(record.ID),
			strings.TrimSpace(record.Name),
			strings.ToLower(strings.TrimSpace(record.Email)),
			fmt.Sprintf("%.2f", record.Score),
			valid,
		}

		if err := writer.Write(outputRow); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

func validateRow(row []string) (Record, error) {
	if len(row) < 4 {
		return Record{}, fmt.Errorf("insufficient columns")
	}

	id, err := strconv.Atoi(row[0])
	if err != nil || id <= 0 {
		return Record{}, fmt.Errorf("invalid ID")
	}

	name := strings.TrimSpace(row[1])
	if name == "" {
		return Record{}, fmt.Errorf("empty name")
	}

	email := strings.TrimSpace(row[2])
	if !strings.Contains(email, "@") {
		return Record{}, fmt.Errorf("invalid email")
	}

	score, err := strconv.ParseFloat(row[3], 64)
	if err != nil || score < 0 || score > 100 {
		return Record{}, fmt.Errorf("invalid score")
	}

	return Record{
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

	fmt.Printf("Data cleaning completed. Output saved to %s\n", outputFile)
}