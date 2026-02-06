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
			return fmt.Errorf("error reading CSV: %w", err)
		}

		cleaned := make([]string, len(record))
		for i, field := range record {
			cleaned[i] = strings.TrimSpace(field)
		}

		if err := writer.Write(cleaned); err != nil {
			return fmt.Errorf("error writing CSV: %w", err)
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

	fmt.Printf("Successfully cleaned %s to %s\n", inputFile, outputFile)
}
package main

import "fmt"

func removeDuplicates(input []int) []int {
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
	data := []int{4, 2, 8, 2, 4, 9, 8, 1}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

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

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	var unique []string
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !dc.seen[trimmed] {
			dc.seen[trimmed] = true
			unique = append(unique, trimmed)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"  apple  ", "banana", "apple", "", "cherry", "banana"}
	unique := cleaner.RemoveDuplicates(data)
	fmt.Println("Deduplicated:", unique)
	
	emails := []string{"test@example.com", "invalid-email", "another@test.org"}
	for _, email := range emails {
		fmt.Printf("Email %s valid: %v\n", email, cleaner.ValidateEmail(email))
	}
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	records []string
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]string, 0),
	}
}

func (dc *DataCleaner) AddRecord(record string) {
	dc.records = append(dc.records, strings.TrimSpace(record))
}

func (dc *DataCleaner) RemoveDuplicates() []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, record := range dc.records {
		if !seen[record] {
			seen[record] = true
			result = append(result, record)
		}
	}

	dc.records = result
	return result
}

func (dc *DataCleaner) ValidateRecords() (valid []string, invalid []string) {
	valid = make([]string, 0)
	invalid = make([]string, 0)

	for _, record := range dc.records {
		if len(record) > 0 && !strings.ContainsAny(record, "!@#$%") {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}

	return valid, invalid
}

func (dc *DataCleaner) GetRecordCount() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()
	
	cleaner.AddRecord("apple")
	cleaner.AddRecord("banana")
	cleaner.AddRecord("apple")
	cleaner.AddRecord("  orange  ")
	cleaner.AddRecord("")
	cleaner.AddRecord("cherry!")
	
	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())
	
	cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", cleaner.GetRecordCount())
	
	valid, invalid := cleaner.ValidateRecords()
	fmt.Printf("Valid records: %v\n", valid)
	fmt.Printf("Invalid records: %v\n", invalid)
}