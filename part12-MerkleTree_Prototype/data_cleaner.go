
package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}
package main

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

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record: %w", err)
		}

		cleaned := make([]string, len(record))
		for i, field := range record {
			cleaned[i] = strings.TrimSpace(field)
		}

		if err := writer.Write(cleaned); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
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

	fmt.Printf("Successfully cleaned data from %s to %s\n", inputFile, outputFile)
}
package utils

import "strings"

func SanitizeInput(input string) string {
    trimmed := strings.TrimSpace(input)
    return strings.Join(strings.Fields(trimmed), " ")
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
	input := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(input)
	fmt.Printf("Original: %v\n", input)
	fmt.Printf("Cleaned: %v\n", cleaned)
}
package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

import (
	"strings"
)

func RemoveDuplicates(lines []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func CleanData(input []string) []string {
	cleaned := make([]string, 0, len(input))
	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return RemoveDuplicates(cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	result := []DataRecord{}
	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			result = append(result, record)
		}
	}
	return result
}

func ValidateRecords(records []DataRecord) []DataRecord {
	for i := range records {
		records[i].Valid = records[i].ID > 0 &&
			len(strings.TrimSpace(records[i].Name)) > 0 &&
			strings.Contains(records[i].Email, "@")
	}
	return records
}

func CleanData(records []DataRecord) []DataRecord {
	unique := RemoveDuplicates(records)
	validated := ValidateRecords(unique)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{0, "Invalid User", "invalid-email", false},
	}

	cleaned := CleanData(sampleData)
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Valid: %v\n", record.ID, record.Name, record.Valid)
	}
}
package main

import (
	"fmt"
	"strings"
)

func CleanData(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	data := []string{" apple ", "banana", " apple", "banana ", "", "  cherry  "}
	cleaned := CleanData(data)
	fmt.Println("Cleaned data:", cleaned)
}