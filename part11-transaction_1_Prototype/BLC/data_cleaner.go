
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

    seen := make(map[string]bool)
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read CSV record: %w", err)
        }

        for i, field := range record {
            record[i] = strings.TrimSpace(field)
        }

        key := strings.Join(record, "|")
        if seen[key] {
            continue
        }
        seen[key] = true

        if err := writer.Write(record); err != nil {
            return fmt.Errorf("failed to write CSV record: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    if err := cleanCSV(os.Args[1], os.Args[2]); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("CSV cleaning completed successfully")
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

func deduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if !seen[email] {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func cleanData(records []DataRecord) []DataRecord {
	emailMap := make(map[string]bool)
	var cleaned []DataRecord

	for _, record := range records {
		record.Email = strings.ToLower(strings.TrimSpace(record.Email))
		record.Valid = validateEmail(record.Email)

		if !emailMap[record.Email] && record.Valid {
			emailMap[record.Email] = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
	}

	cleaned := cleanData(records)
	fmt.Printf("Cleaned records: %d\n", len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
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