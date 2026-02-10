
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
}package datautils

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
}package utils

import "strings"

func SanitizeInput(input string) string {
    trimmed := strings.TrimSpace(input)
    return strings.Join(strings.Fields(trimmed), " ")
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
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package datautils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(input)

	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized := spaceRegex.ReplaceAllString(trimmed, " ")

	// Remove non-printable characters
	var result strings.Builder
	for _, r := range normalized {
		if unicode.IsPrint(r) {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func NormalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
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

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	header = sanitizeHeaders(header)
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		cleaned := make([]string, len(record))
		for i, field := range record {
			cleaned[i] = strings.TrimSpace(field)
		}

		if isEmptyRecord(cleaned) {
			continue
		}

		if err := writer.Write(cleaned); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

func sanitizeHeaders(headers []string) []string {
	sanitized := make([]string, len(headers))
	for i, h := range headers {
		sanitized[i] = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(h), " ", "_"))
	}
	return sanitized
}

func isEmptyRecord(fields []string) bool {
	for _, f := range fields {
		if f != "" {
			return false
		}
	}
	return true
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

	fmt.Printf("Successfully cleaned data. Output saved to %s\n", outputFile)
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

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(normalized) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && len(item) < 100
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"apple", "Apple", "banana", "  BANANA  ", "", "cherry", "cherry"}
	
	fmt.Println("Original data:", data)
	
	cleaned := cleaner.RemoveDuplicates(data)
	fmt.Println("Cleaned data:", cleaned)
	
	cleaner.Reset()
	
	moreData := []string{"grape", "GRAPE", "kiwi"}
	cleaned2 := cleaner.RemoveDuplicates(moreData)
	fmt.Println("Second batch:", cleaned2)
}
package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput cleans user-provided strings by:
// - Removing leading/trailing whitespace
// - Collapsing multiple spaces into single spaces
// - Removing non-alphanumeric characters (configurable)
// - Converting to lowercase (optional)
func SanitizeInput(input string, allowSpecialChars bool, toLower bool) string {
	if input == "" {
		return input
	}

	// Trim whitespace
	cleaned := strings.TrimSpace(input)

	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned = spaceRegex.ReplaceAllString(cleaned, " ")

	// Remove non-alphanumeric characters if not allowed
	if !allowSpecialChars {
		specialCharRegex := regexp.MustCompile(`[^a-zA-Z0-9\s\-_]`)
		cleaned = specialCharRegex.ReplaceAllString(cleaned, "")
	}

	// Convert to lowercase if requested
	if toLower {
		cleaned = strings.ToLower(cleaned)
	}

	return cleaned
}

// ValidateEmail checks if a string is a valid email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
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

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return len(email) > 5 && len(email) < 255
}

func validateRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		record.Valid = validateEmail(record.Email)
		if record.Valid {
			valid = append(valid, record)
		}
	}
	return valid
}

func processData(records []DataRecord) []DataRecord {
	unique := deduplicateRecords(records)
	valid := validateRecords(unique)
	return valid
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob Wilson", "invalid-email", false},
	}

	cleaned := processData(sampleData)
	fmt.Printf("Processed %d records, %d valid after cleaning\n", len(sampleData), len(cleaned))
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}