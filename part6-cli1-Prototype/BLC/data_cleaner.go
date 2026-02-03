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

	headerProcessed := false
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}

		if !headerProcessed {
			headerProcessed = true
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing header: %w", err)
			}
			continue
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedField := strings.TrimSpace(field)
			cleanedField = strings.ToLower(cleanedField)
			cleanedRecord[i] = cleanedField
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("error writing record: %w", err)
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

	fmt.Printf("Successfully cleaned data. Output saved to %s\n", outputFile)
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

	strings := []string{"apple", "banana", "apple", "orange", "banana"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID        int
	Name      string
	Email     string
	Age       int
	Active    bool
	Timestamp string
}

func parseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if lineNum == 0 {
			lineNum++
			continue
		}

		record, err := parseRecord(line)
		if err != nil {
			return nil, fmt.Errorf("parse error at line %d: %w", lineNum, err)
		}

		records = append(records, record)
		lineNum++
	}

	return records, nil
}

func parseRecord(fields []string) (Record, error) {
	if len(fields) != 6 {
		return Record{}, fmt.Errorf("invalid field count: %d", len(fields))
	}

	id, err := strconv.Atoi(fields[0])
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID: %v", err)
	}

	name := strings.TrimSpace(fields[1])
	if name == "" {
		return Record{}, fmt.Errorf("name cannot be empty")
	}

	email := strings.TrimSpace(fields[2])
	if !strings.Contains(email, "@") {
		return Record{}, fmt.Errorf("invalid email format")
	}

	age, err := strconv.Atoi(fields[3])
	if err != nil || age < 0 || age > 150 {
		return Record{}, fmt.Errorf("invalid age: %v", err)
	}

	active := false
	if fields[4] == "true" {
		active = true
	} else if fields[4] != "false" {
		return Record{}, fmt.Errorf("invalid active flag: %s", fields[4])
	}

	timestamp := strings.TrimSpace(fields[5])
	if timestamp == "" {
		return Record{}, fmt.Errorf("timestamp cannot be empty")
	}

	return Record{
		ID:        id,
		Name:      name,
		Email:     email,
		Age:       age,
		Active:    active,
		Timestamp: timestamp,
	}, nil
}

func validateRecords(records []Record) ([]Record, []string) {
	validRecords := []Record{}
	errors := []string{}

	for _, record := range records {
		if record.Age < 18 {
			errors = append(errors, fmt.Sprintf("record ID %d: age must be 18 or older", record.ID))
			continue
		}

		if !record.Active {
			errors = append(errors, fmt.Sprintf("record ID %d: user is inactive", record.ID))
			continue
		}

		validRecords = append(validRecords, record)
	}

	return validRecords, errors
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_cleaner <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := parseCSVFile(filename)
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed %d records from %s\n", len(records), filename)

	validRecords, errors := validateRecords(records)

	fmt.Printf("\nValidation Results:\n")
	fmt.Printf("Valid records: %d\n", len(validRecords))
	fmt.Printf("Validation errors: %d\n", len(errors))

	if len(errors) > 0 {
		fmt.Println("\nErrors:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
	}

	if len(validRecords) > 0 {
		fmt.Println("\nValid Records:")
		for _, record := range validRecords {
			fmt.Printf("  ID: %d, Name: %s, Email: %s, Age: %d\n",
				record.ID, record.Name, record.Email, record.Age)
		}
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}