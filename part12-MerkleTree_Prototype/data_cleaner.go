
package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(input)
	
	// Remove extra internal whitespace
	re := regexp.MustCompile(`\s+`)
	normalized := re.ReplaceAllString(trimmed, " ")
	
	// Convert to lowercase
	lowercased := strings.ToLower(normalized)
	
	return lowercased
}

func RemoveSpecialChars(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, input)
}

func NormalizeWhitespace(input string) string {
	// Replace various whitespace characters with standard space
	re := regexp.MustCompile(`[\t\n\r\f\v]+`)
	return re.ReplaceAllString(input, " ")
}package main

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

func parseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataRecord
	for i, row := range records {
		if i == 0 {
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
		email := strings.TrimSpace(row[2])
		score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
		if err != nil {
			continue
		}

		record := DataRecord{
			ID:    id,
			Name:  name,
			Email: email,
			Score: score,
		}
		data = append(data, record)
	}

	return data, nil
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func cleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	for _, record := range records {
		if record.ID <= 0 {
			continue
		}
		if record.Name == "" {
			continue
		}
		if !validateEmail(record.Email) {
			continue
		}
		if record.Score < 0 || record.Score > 100 {
			continue
		}
		cleaned = append(cleaned, record)
	}
	return cleaned
}

func writeCleanCSV(filename string, records []DataRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Name", "Email", "Score"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range records {
		row := []string{
			strconv.Itoa(record.ID),
			record.Name,
			record.Email,
			strconv.FormatFloat(record.Score, 'f', 2, 64),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	records, err := parseCSVFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	fmt.Printf("Read %d records from %s\n", len(records), inputFile)

	cleaned := cleanData(records)
	fmt.Printf("Cleaned data: %d valid records\n", len(cleaned))

	if err := writeCleanCSV(outputFile, cleaned); err != nil {
		fmt.Printf("Error writing CSV: %v\n", err)
		return
	}

	fmt.Printf("Cleaned data written to %s\n", outputFile)
}