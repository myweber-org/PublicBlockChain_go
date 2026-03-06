
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
}