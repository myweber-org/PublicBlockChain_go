
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
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func validateRecords(records []DataRecord) []DataRecord {
	var validated []DataRecord
	for _, record := range records {
		record.Valid = validateEmail(record.Email)
		validated = append(validated, record)
	}
	return validated
}

func cleanData(records []DataRecord) []DataRecord {
	unique := deduplicateRecords(records)
	validated := validateRecords(unique)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob Wilson", "invalid-email", false},
	}

	cleaned := cleanData(sampleData)
	
	for _, record := range cleaned {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Name: %s, Status: %s\n", record.ID, record.Name, status)
	}
}
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

type DataRecord struct {
	ID      string
	Name    string
	Email   string
	Valid   bool
	Dupe    bool
}

func readCSV(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	header := true

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if header {
			header = false
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Name:  strings.TrimSpace(row[1]),
			Email: strings.TrimSpace(row[2]),
		}
		records = append(records, record)
	}
	return records, nil
}

func validateRecords(records []DataRecord) []DataRecord {
	for i := range records {
		records[i].Valid = len(records[i].Email) > 0 && strings.Contains(records[i].Email, "@")
	}
	return records
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	for i := range records {
		key := records[i].Email
		if seen[key] {
			records[i].Dupe = true
		} else {
			seen[key] = true
			records[i].Dupe = false
		}
	}
	return records
}

func writeCleanCSV(filename string, records []DataRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Name", "Email", "Valid", "Duplicate"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range records {
		if record.Valid && !record.Dupe {
			row := []string{
				record.ID,
				record.Name,
				record.Email,
				fmt.Sprintf("%t", record.Valid),
				fmt.Sprintf("%t", record.Dupe),
			}
			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run data_cleaner.go <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	records, err := readCSV(inputFile)
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		os.Exit(1)
	}

	records = validateRecords(records)
	records = deduplicateRecords(records)

	err = writeCleanCSV(outputFile, records)
	if err != nil {
		fmt.Printf("Error writing CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Cleaned data written to %s\n", outputFile)
}