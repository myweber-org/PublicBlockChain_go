
package main

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	return input
}

func TransformToSlug(text string) string {
	text = strings.ToLower(text)
	text = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(text, "-")
	text = strings.Trim(text, "-")
	return text
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
	ID        int
	Name      string
	Value     float64
	Timestamp string
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNum := 0

	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNum, len(row))
		}

		record, err := validateAndCreateRecord(row, lineNum)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func validateAndCreateRecord(row []string, lineNum int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return record, fmt.Errorf("empty name at line %d", lineNum)
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
	}
	record.Value = value

	timestamp := strings.TrimSpace(row[3])
	if timestamp == "" {
		return record, fmt.Errorf("empty timestamp at line %d", lineNum)
	}
	record.Timestamp = timestamp

	return record, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, float64) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var min, max float64

	for i, record := range records {
		sum += record.Value
		if i == 0 {
			min = record.Value
			max = record.Value
		} else {
			if record.Value < min {
				min = record.Value
			}
			if record.Value > max {
				max = record.Value
			}
		}
	}

	average := sum / float64(len(records))
	return average, min, max
}

func FilterRecords(records []DataRecord, minValue, maxValue float64) []DataRecord {
	filtered := []DataRecord{}
	for _, record := range records {
		if record.Value >= minValue && record.Value <= maxValue {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := ParseCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully parsed %d records\n", len(records))

	avg, min, max := CalculateStatistics(records)
	fmt.Printf("Statistics - Average: %.2f, Min: %.2f, Max: %.2f\n", avg, min, max)

	filtered := FilterRecords(records, 10.0, 100.0)
	fmt.Printf("Records between 10.0 and 100.0: %d\n", len(filtered))
}