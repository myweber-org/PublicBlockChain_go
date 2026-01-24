
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
    fmt.Println("Original:", data)
    fmt.Println("Cleaned:", cleaned)
}
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

func parseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
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

		if !isValidEmail(email) {
			continue
		}

		data = append(data, DataRecord{
			ID:    id,
			Name:  name,
			Email: email,
			Score: score,
		})
	}

	return data, nil
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func calculateAverageScore(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0.0
	}

	var total float64
	for _, record := range records {
		total += record.Score
	}
	return total / float64(len(records))
}

func filterHighScorers(records []DataRecord, threshold float64) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Score >= threshold {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func writeCleanData(filename string, records []DataRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Name", "Email", "Score"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, record := range records {
		row := []string{
			strconv.Itoa(record.ID),
			record.Name,
			record.Email,
			strconv.FormatFloat(record.Score, 'f', 2, 64),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_cleaner <input_file.csv>")
		return
	}

	inputFile := os.Args[1]
	records, err := parseCSVFile(inputFile)
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		return
	}

	fmt.Printf("Parsed %d valid records\n", len(records))
	
	averageScore := calculateAverageScore(records)
	fmt.Printf("Average score: %.2f\n", averageScore)

	highScorers := filterHighScorers(records, 80.0)
	fmt.Printf("Found %d high scorers (>= 80.0)\n", len(highScorers))

	outputFile := "cleaned_data.csv"
	if err := writeCleanData(outputFile, records); err != nil {
		fmt.Printf("Error writing cleaned data: %v\n", err)
		return
	}

	fmt.Printf("Cleaned data written to %s\n", outputFile)
}