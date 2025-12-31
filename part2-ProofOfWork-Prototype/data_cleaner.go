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

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	trimmedHeaders := make([]string, len(headers))
	for i, h := range headers {
		trimmedHeaders[i] = strings.TrimSpace(h)
	}
	if err := writer.Write(trimmedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
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
			if cleaned[i] == "" {
				cleaned[i] = "N/A"
			}
		}
		if err := writer.Write(cleaned); err != nil {
			return fmt.Errorf("failed to write cleaned record: %w", err)
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

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
    "fmt"
    "strings"
)

type DataRecord struct {
    ID    int
    Email string
    Phone string
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
    seen := make(map[string]bool)
    unique := []DataRecord{}
    for _, record := range records {
        key := fmt.Sprintf("%s|%s", record.Email, record.Phone)
        if !seen[key] {
            seen[key] = true
            unique = append(unique, record)
        }
    }
    return unique
}

func ValidateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func ValidatePhone(phone string) bool {
    return len(phone) >= 10 && strings.HasPrefix(phone, "+")
}

func CleanData(records []DataRecord) []DataRecord {
    validRecords := []DataRecord{}
    for _, record := range records {
        if ValidateEmail(record.Email) && ValidatePhone(record.Phone) {
            validRecords = append(validRecords, record)
        }
    }
    return DeduplicateRecords(validRecords)
}

func main() {
    sampleData := []DataRecord{
        {1, "test@example.com", "+1234567890"},
        {2, "invalid-email", "+0987654321"},
        {3, "test@example.com", "+1234567890"},
        {4, "another@test.org", "+445551234"},
    }

    cleaned := CleanData(sampleData)
    fmt.Printf("Original count: %d\n", len(sampleData))
    fmt.Printf("Cleaned count: %d\n", len(cleaned))
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", record.ID, record.Email, record.Phone)
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
    Phone string
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
    seen := make(map[int]bool)
    result := []DataRecord{}
    
    for _, record := range records {
        if !seen[record.ID] {
            seen[record.ID] = true
            result = append(result, record)
        }
    }
    return result
}

func ValidateEmail(email string) bool {
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    return len(parts[0]) > 0 && len(parts[1]) > 0
}

func CleanPhoneNumber(phone string) string {
    var builder strings.Builder
    for _, ch := range phone {
        if ch >= '0' && ch <= '9' {
            builder.WriteRune(ch)
        }
    }
    return builder.String()
}

func ProcessRecords(records []DataRecord) []DataRecord {
    uniqueRecords := RemoveDuplicates(records)
    
    for i := range uniqueRecords {
        uniqueRecords[i].Phone = CleanPhoneNumber(uniqueRecords[i].Phone)
    }
    
    validRecords := []DataRecord{}
    for _, record := range uniqueRecords {
        if ValidateEmail(record.Email) {
            validRecords = append(validRecords, record)
        }
    }
    
    return validRecords
}

func main() {
    sampleData := []DataRecord{
        {1, "user@example.com", "(123) 456-7890"},
        {2, "invalid-email", "555-1234"},
        {1, "user@example.com", "1234567890"},
        {3, "another@test.org", "+1-800-555-0199"},
    }
    
    cleaned := ProcessRecords(sampleData)
    
    fmt.Printf("Processed %d records\n", len(cleaned))
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", 
            record.ID, record.Email, record.Phone)
    }
}