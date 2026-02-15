
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

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"apple", "Apple ", "banana", "BANANA", "  cherry  "}
	
	fmt.Println("Original data:", data)
	
	deduped := cleaner.Deduplicate(data)
	fmt.Println("Deduplicated:", deduped)
	
	cleaner.Reset()
	
	testValue := "TestValue"
	fmt.Printf("Is '%s' duplicate? %v\n", testValue, cleaner.IsDuplicate(testValue))
	fmt.Printf("Is '%s' duplicate? %v\n", strings.ToLower(testValue), cleaner.IsDuplicate(strings.ToLower(testValue)))
}
package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Score int
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if !strings.Contains(record.Email, "@") {
		return errors.New("invalid email format")
	}
	if record.Score < 0 || record.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}
	return nil
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func CleanData(records []DataRecord) ([]DataRecord, error) {
	var cleaned []DataRecord

	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			fmt.Printf("Skipping invalid record %s: %v\n", record.ID, err)
			continue
		}
		cleaned = append(cleaned, record)
	}

	cleaned = DeduplicateRecords(cleaned)
	return cleaned, nil
}

func main() {
	records := []DataRecord{
		{"A1", "test@example.com", 85},
		{"A2", "invalid-email", 92},
		{"A1", "duplicate@example.com", 78},
		{"A3", "another@test.com", 105},
		{"A4", "valid@domain.com", 67},
	}

	cleaned, err := CleanData(records)
	if err != nil {
		fmt.Printf("Error cleaning data: %v\n", err)
		return
	}

	fmt.Printf("Original records: %d\n", len(records))
	fmt.Printf("Cleaned records: %d\n", len(cleaned))
	for _, record := range cleaned {
		fmt.Printf("ID: %s, Email: %s, Score: %d\n", record.ID, record.Email, record.Score)
	}
}package main

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
	data := []int{5, 2, 8, 2, 5, 9, 8, 1}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package utils

func RemoveDuplicates(nums []int) []int {
    if len(nums) == 0 {
        return nums
    }
    
    seen := make(map[int]bool)
    result := make([]int, 0, len(nums))
    
    for _, num := range nums {
        if !seen[num] {
            seen[num] = true
            result = append(result, num)
        }
    }
    
    return result
}package main

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
	data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6, 7}
	cleanedData := removeDuplicates(data)
	fmt.Println("Original data:", data)
	fmt.Println("Cleaned data:", cleanedData)
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

func readCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []DataRecord
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		if len(line) != 4 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(line[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(line[1])
		email := strings.TrimSpace(line[2])
		score, err := strconv.ParseFloat(strings.TrimSpace(line[3]), 64)
		if err != nil {
			continue
		}

		record := DataRecord{
			ID:    id,
			Name:  name,
			Email: email,
			Score: score,
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
	return true
}

func cleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	for _, record := range records {
		if record.ID <= 0 {
			continue
		}
		if len(record.Name) == 0 {
			continue
		}
		if !validateEmail(record.Email) {
			continue
		}
		if record.Score < 0 || record.Score > 100 {
			continue
		}
		cleaned = append(cleaned, record)
	}
	return cleaned
}

func calculateAverageScore(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}
	var total float64
	for _, record := range records {
		total += record.Score
	}
	return total / float64(len(records))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run data_cleaner.go <csv_file>")
		return
	}

	filename := os.Args[1]
	records, err := readCSVFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Printf("Read %d records from %s\n", len(records), filename)

	cleaned := cleanData(records)
	fmt.Printf("After cleaning: %d valid records\n", len(cleaned))

	average := calculateAverageScore(cleaned)
	fmt.Printf("Average score: %.2f\n", average)

	for i, record := range cleaned {
		if i < 5 {
			fmt.Printf("Record %d: ID=%d, Name=%s, Email=%s, Score=%.1f\n",
				i+1, record.ID, record.Name, record.Email, record.Score)
		}
	}
}package main

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