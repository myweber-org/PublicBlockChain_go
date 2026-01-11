
package main

import "fmt"

func removeDuplicates(nums []int) []int {
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
	data := []int{4, 2, 7, 2, 4, 9, 7, 1}
	cleaned := removeDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct{}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func (dc *DataCleaner) TrimWhitespace(items []string) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func main() {
	cleaner := &DataCleaner{}
	data := []string{"  apple ", "banana", "  apple ", " cherry", "banana "}

	trimmed := cleaner.TrimWhitespace(data)
	unique := cleaner.RemoveDuplicates(trimmed)

	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", unique)
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Name  string
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) {
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
	seen := make(map[string]bool)
	result := make([]DataRecord, 0)

	for _, record := range dc.records {
		key := fmt.Sprintf("%d|%s|%s", record.ID, strings.ToLower(record.Email), strings.ToLower(record.Name))
		if !seen[key] {
			seen[key] = true
			result = append(result, record)
		}
	}

	dc.records = result
	return result
}

func (dc *DataCleaner) ValidateEmails() (valid []DataRecord, invalid []DataRecord) {
	for _, record := range dc.records {
		if strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".") {
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

	cleaner.AddRecord(DataRecord{ID: 1, Email: "user@example.com", Name: "John Doe"})
	cleaner.AddRecord(DataRecord{ID: 2, Email: "user@example.com", Name: "John Doe"})
	cleaner.AddRecord(DataRecord{ID: 3, Email: "jane@test.org", Name: "Jane Smith"})
	cleaner.AddRecord(DataRecord{ID: 4, Email: "invalid-email", Name: "Bad Data"})

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

	cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", cleaner.GetRecordCount())

	valid, invalid := cleaner.ValidateEmails()
	fmt.Printf("Valid emails: %d, Invalid emails: %d\n", len(valid), len(invalid))
}package datautils

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
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type Cleaner struct {
	TrimSpaces bool
	RemoveEmpty bool
}

func NewCleaner() *Cleaner {
	return &Cleaner{
		TrimSpaces: true,
		RemoveEmpty: true,
	}
}

func (c *Cleaner) ProcessRow(row []string) []string {
	var result []string
	for _, field := range row {
		processed := field
		if c.TrimSpaces {
			processed = strings.TrimSpace(processed)
		}
		if !c.RemoveEmpty || processed != "" {
			result = append(result, processed)
		}
	}
	return result
}

func (c *Cleaner) CleanCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
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
			return fmt.Errorf("read csv: %w", err)
		}

		cleaned := c.ProcessRow(record)
		if len(cleaned) > 0 {
			if err := writer.Write(cleaned); err != nil {
				return fmt.Errorf("write csv: %w", err)
			}
		}
	}
	return nil
}

func main() {
	cleaner := NewCleaner()
	if err := cleaner.CleanCSV("input.csv", "output.csv"); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("CSV cleaning completed successfully")
}