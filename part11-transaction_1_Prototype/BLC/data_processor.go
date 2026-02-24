
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
package main

import (
	"errors"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Category  string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.After(time.Now()) {
		return errors.New("timestamp cannot be in the future")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	record.Category = strings.ToUpper(strings.TrimSpace(record.Category))
	record.Value = roundToTwoDecimals(record.Value)
	return record
}

func roundToTwoDecimals(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, err
		}
		processed = append(processed, TransformRecord(record))
	}
	return processed, nil
}
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
		}

		record, err := parseRow(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		return nil, errors.New("no valid records found in file")
	}

	return records, nil
}

func parseRow(row []string, lineNumber int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return record, fmt.Errorf("empty name at line %d", lineNumber)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
	}
	record.Active = active

	return record, nil
}

func ValidateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID %d found", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("invalid value %f for ID %d: cannot be negative", record.Value, record.ID)
		}
	}

	return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var activeCount int
	var maxValue float64

	for _, record := range records {
		sum += record.Value
		if record.Value > maxValue {
			maxValue = record.Value
		}
		if record.Active {
			activeCount++
		}
	}

	average := sum / float64(len(records))
	return average, maxValue, activeCount
}package main

import "fmt"

func calculateAverage(numbers []int) float64 {
    if len(numbers) == 0 {
        return 0
    }
    sum := 0
    for _, num := range numbers {
        sum += num
    }
    return float64(sum) / float64(len(numbers))
}

func main() {
    data := []int{10, 20, 30, 40, 50}
    avg := calculateAverage(data)
    fmt.Printf("Average: %.2f\n", avg)
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

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}

		cleaned := make([]string, len(record))
		for i, field := range record {
			cleaned[i] = strings.TrimSpace(field)
		}

		if err := writer.Write(cleaned); err != nil {
			return fmt.Errorf("error writing CSV: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	err := cleanCSV(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("CSV processing completed successfully")
}