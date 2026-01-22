
package main

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
    Validated bool
}

func cleanEmail(email string) string {
    return strings.ToLower(strings.TrimSpace(email))
}

func validateAge(age int) bool {
    return age >= 18 && age <= 120
}

func processCSVFile(inputPath string) ([]Record, error) {
    file, err := os.Open(inputPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
    lineNumber := 0

    for {
        row, err := reader.Read()
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

        if len(row) < 4 {
            continue
        }

        id, err := strconv.Atoi(strings.TrimSpace(row[0]))
        if err != nil {
            continue
        }

        name := strings.TrimSpace(row[1])
        email := cleanEmail(row[2])
        age, err := strconv.Atoi(strings.TrimSpace(row[3]))
        if err != nil {
            continue
        }

        record := Record{
            ID:        id,
            Name:      name,
            Email:     email,
            Age:       age,
            Validated: validateAge(age),
        }

        records = append(records, record)
    }

    return records, nil
}

func generateReport(records []Record) {
    validCount := 0
    for _, record := range records {
        if record.Validated {
            validCount++
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Valid records: %d\n", validCount)
    fmt.Printf("Invalid records: %d\n", len(records)-validCount)
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run data_cleaner.go <input_file.csv>")
        return
    }

    records, err := processCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        return
    }

    generateReport(records)
}package main

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
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	for i := range headers {
		headers[i] = strings.TrimSpace(headers[i])
	}

	if err := writer.Write(headers); err != nil {
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

		for i := range record {
			record[i] = strings.TrimSpace(record[i])
		}

		key := strings.Join(record, "|")
		if seen[key] {
			continue
		}
		seen[key] = true

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
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

	fmt.Printf("Successfully cleaned data written to %s\n", outputFile)
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
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
	data := []int{4, 2, 8, 2, 4, 9, 8, 1}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}