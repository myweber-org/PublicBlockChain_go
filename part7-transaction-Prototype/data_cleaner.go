package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func RemoveDuplicateRows(inputPath, outputPath string) error {
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

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}

		key := fmt.Sprintf("%v", record)
		if !seen[key] {
			seen[key] = true
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing CSV: %w", err)
			}
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	err := RemoveDuplicateRows(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Duplicate removal completed successfully")
}
package main

import (
	"fmt"
	"strings"
)

func CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	return lower
}

func RemoveDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		cleaned := CleanString(elements[v])
		if !encountered[cleaned] {
			encountered[cleaned] = true
			result = append(result, cleaned)
		}
	}
	return result
}

func main() {
	data := []string{" Apple", "banana ", "apple", " Banana", "Cherry", "cherry "}
	unique := RemoveDuplicates(data)
	fmt.Println("Cleaned unique data:", unique)
}package main

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
	numbers := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Cleaned:", uniqueNumbers)
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6, 7}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedRecords map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processedRecords: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(records []string) []string {
	uniqueRecords := []string{}
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] {
			dc.processedRecords[normalized] = true
			uniqueRecords = append(uniqueRecords, record)
		}
	}
	return uniqueRecords
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (dc *DataCleaner) Reset() {
	dc.processedRecords = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	emails := []string{
		"user@example.com",
		"USER@example.com",
		"test@domain.org",
		"user@example.com",
		"invalid-email",
	}
	
	uniqueEmails := cleaner.RemoveDuplicates(emails)
	fmt.Printf("Original: %v\n", emails)
	fmt.Printf("Deduplicated: %v\n", uniqueEmails)
	
	for _, email := range uniqueEmails {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("Valid: %s\n", email)
		} else {
			fmt.Printf("Invalid: %s\n", email)
		}
	}
	
	cleaner.Reset()
	fmt.Println("Cleaner has been reset")
}