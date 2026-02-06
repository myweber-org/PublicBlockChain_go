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
	Valid bool
}

func parseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNum++
		if lineNum == 1 {
			continue
		}

		if len(line) != 5 {
			continue
		}

		id, idErr := strconv.Atoi(strings.TrimSpace(line[0]))
		name := strings.TrimSpace(line[1])
		email := strings.TrimSpace(line[2])
		score, scoreErr := strconv.ParseFloat(strings.TrimSpace(line[3]), 64)
		valid := strings.ToLower(strings.TrimSpace(line[4])) == "true"

		if idErr != nil || scoreErr != nil {
			continue
		}

		if !validateEmail(email) {
			valid = false
		}

		record := DataRecord{
			ID:    id,
			Name:  name,
			Email: email,
			Score: score,
			Valid: valid,
		}
		records = append(records, record)
	}

	return records, nil
}

func validateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func filterValidRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func calculateAverageScore(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0.0
	}

	total := 0.0
	for _, record := range records {
		total += record.Score
	}
	return total / float64(len(records))
}

func generateReport(records []DataRecord) {
	validRecords := filterValidRecords(records)
	invalidCount := len(records) - len(validRecords)
	averageScore := calculateAverageScore(validRecords)

	fmt.Printf("Data Cleaning Report\n")
	fmt.Printf("====================\n")
	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Valid records: %d\n", len(validRecords))
	fmt.Printf("Invalid records: %d\n", invalidCount)
	fmt.Printf("Average score (valid records): %.2f\n", averageScore)

	if len(validRecords) > 0 {
		fmt.Printf("\nTop 5 valid records by score:\n")
		for i := 0; i < len(validRecords) && i < 5; i++ {
			record := validRecords[i]
			fmt.Printf("%d. %s (Score: %.1f, Email: %s)\n",
				i+1, record.Name, record.Score, record.Email)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run data_cleaner.go <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := parseCSVFile(filename)
	if err != nil {
		fmt.Printf("Error reading CSV file: %v\n", err)
		os.Exit(1)
	}

	generateReport(records)
}