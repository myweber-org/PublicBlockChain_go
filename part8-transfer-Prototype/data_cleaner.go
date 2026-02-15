
package datautils

func DeduplicateInts(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, v := range input {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func FilterPositiveInts(input []int) []int {
	result := []int{}
	for _, v := range input {
		if v > 0 {
			result = append(result, v)
		}
	}
	return result
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

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

	seen := make(map[string]bool)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	for i := range headers {
		headers[i] = strings.TrimSpace(headers[i])
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		for i := range record {
			record[i] = strings.TrimSpace(record[i])
		}

		key := strings.Join(record, "|")
		if seen[key] {
			continue
		}
		seen[key] = true

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
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

	if err := cleanCSV(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully cleaned data written to %s\n", outputFile)
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

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Score float64
}

func cleanCSVData(inputPath, outputPath string) error {
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

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	validCount := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read row: %w", err)
		}

		recordCount++
		cleanedRow, isValid := processRow(row)

		if isValid {
			if err := writer.Write(cleanedRow); err != nil {
				return fmt.Errorf("failed to write row: %w", err)
			}
			validCount++
		}
	}

	fmt.Printf("Processed %d records, %d valid records written to %s\n", recordCount, validCount, outputPath)
	return nil
}

func processRow(row []string) ([]string, bool) {
	if len(row) != 4 {
		return nil, false
	}

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil || id <= 0 {
		return nil, false
	}

	name := strings.TrimSpace(row[1])
	if name == "" || len(name) > 100 {
		return nil, false
	}

	email := strings.TrimSpace(row[2])
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return nil, false
	}

	score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
	if err != nil || score < 0 || score > 100 {
		return nil, false
	}

	return []string{
		strconv.Itoa(id),
		name,
		strings.ToLower(email),
		fmt.Sprintf("%.2f", score),
	}, true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run data_cleaner.go <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSVData(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	input := []int{1, 2, 2, 3, 4, 4, 5, 6, 6}
	cleaned := RemoveDuplicates(input)
	fmt.Printf("Original: %v\n", input)
	fmt.Printf("Cleaned: %v\n", cleaned)
}package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	for _, record := range records {
		if !seen[record.Email] {
			seen[record.Email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	for i := range records {
		records[i].Valid = strings.Contains(records[i].Email, "@") && strings.Contains(records[i].Email, ".")
	}
	return records
}

func CleanData(records []DataRecord) []DataRecord {
	unique := RemoveDuplicates(records)
	validated := ValidateEmails(unique)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "invalid-email", false},
		{3, "user@example.com", false},
		{4, "another@test.org", false},
	}

	cleaned := CleanData(sampleData)
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
	}
}