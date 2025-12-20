package main

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
}package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID    string
    Name  string
    Email string
    Valid bool
}

func cleanString(s string) string {
    return strings.TrimSpace(strings.ToLower(s))
}

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processCSVFile(inputPath string) ([]DataRecord, error) {
    file, err := os.Open(inputPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []DataRecord

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if len(row) < 3 {
            continue
        }

        record := DataRecord{
            ID:    cleanString(row[0]),
            Name:  cleanString(row[1]),
            Email: cleanString(row[2]),
        }
        record.Valid = validateEmail(record.Email)

        records = append(records, record)
    }

    return records, nil
}

func writeCleanData(records []DataRecord, outputPath string) error {
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    header := []string{"id", "name", "email", "valid"}
    if err := writer.Write(header); err != nil {
        return err
    }

    for _, record := range records {
        row := []string{
            record.ID,
            record.Name,
            record.Email,
            fmt.Sprintf("%t", record.Valid),
        }
        if err := writer.Write(row); err != nil {
            return err
        }
    }

    return nil
}

func main() {
    inputFile := "raw_data.csv"
    outputFile := "cleaned_data.csv"

    records, err := processCSVFile(inputFile)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        return
    }

    fmt.Printf("Processed %d records\n", len(records))

    validCount := 0
    for _, record := range records {
        if record.Valid {
            validCount++
        }
    }
    fmt.Printf("Valid records: %d\n", validCount)

    if err := writeCleanData(records, outputFile); err != nil {
        fmt.Printf("Error writing output: %v\n", err)
        return
    }

    fmt.Println("Data cleaning completed successfully")
}