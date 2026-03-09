
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

	csvReader := csv.NewReader(inputFile)
	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	headers, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if err := csvWriter.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	validCount := 0

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		recordCount++

		if len(record) != 4 {
			continue
		}

		cleanRecord := make([]string, 4)

		id, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil || id <= 0 {
			continue
		}
		cleanRecord[0] = strconv.Itoa(id)

		name := strings.TrimSpace(record[1])
		if name == "" {
			continue
		}
		cleanRecord[1] = name

		email := strings.TrimSpace(record[2])
		if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
			continue
		}
		cleanRecord[2] = strings.ToLower(email)

		score, err := strconv.ParseFloat(strings.TrimSpace(record[3]), 64)
		if err != nil || score < 0 || score > 100 {
			continue
		}
		cleanRecord[3] = strconv.FormatFloat(score, 'f', 2, 64)

		if err := csvWriter.Write(cleanRecord); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}

		validCount++
	}

	fmt.Printf("Processed %d records, %d valid records written to %s\n", recordCount, validCount, outputPath)
	return nil
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

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_cleaner <input.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := strings.TrimSuffix(inputFile, ".csv") + "_cleaned.csv"

	inFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	seen := make(map[string]bool)
	headerWritten := false

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading CSV: %v\n", err)
			os.Exit(1)
		}

		if !headerWritten {
			err = writer.Write(record)
			if err != nil {
				fmt.Printf("Error writing header: %v\n", err)
				os.Exit(1)
			}
			headerWritten = true
			continue
		}

		key := strings.Join(record, "|")
		if !seen[key] {
			seen[key] = true
			err = writer.Write(record)
			if err != nil {
				fmt.Printf("Error writing record: %v\n", err)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("Cleaned data written to: %s\n", outputFile)
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
    unique := []string{}
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
    fmt.Println("Unique items:", unique)
    
    emails := []string{"test@example.com", "invalid-email", "user@domain"}
    for _, email := range emails {
        fmt.Printf("Email %s valid: %v\n", email, cleaner.ValidateEmail(email))
    }
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func DeduplicateEmails(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] && email != "" {
			seen[email] = true
			record.Email = email
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmailFormat(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func CleanDataset(records []DataRecord) []DataRecord {
	deduped := DeduplicateEmails(records)
	var cleaned []DataRecord

	for _, record := range deduped {
		record.Valid = ValidateEmailFormat(record.Email)
		if record.Valid {
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain", false},
		{5, "user@example.com", false},
		{6, "new@test.org", false},
	}

	cleaned := CleanDataset(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d valid unique records\n", len(cleaned))

	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s\n", record.ID, record.Email)
	}
}