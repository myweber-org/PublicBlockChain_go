package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	if !validPattern.MatchString(trimmed) {
		return "", false
	}
	return trimmed, true
}

func ValidateEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(email)
}

func ProcessUserData(username, email string) (map[string]interface{}, error) {
	sanitizedUsername, ok := SanitizeUsername(username)
	if !ok {
		return nil, &InvalidDataError{Field: "username", Value: username}
	}

	if !ValidateEmail(email) {
		return nil, &InvalidDataError{Field: "email", Value: email}
	}

	return map[string]interface{}{
		"username": sanitizedUsername,
		"email":    strings.ToLower(email),
		"status":   "processed",
	}, nil
}

type InvalidDataError struct {
	Field string
	Value string
}

func (e *InvalidDataError) Error() string {
	return "invalid data for field: " + e.Field
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
	ID    int
	Name  string
	Value float64
	Valid bool
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

		if len(row) < 4 {
			return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
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

func parseRow(row []string, lineNum int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return DataRecord{}, fmt.Errorf("empty name at line %d", lineNum)
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
	}
	record.Value = value

	valid, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid boolean at line %d: %w", lineNum, err)
	}
	record.Valid = valid

	return record, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0.0
	}

	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	return sum / float64(len(records))
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

type Record struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func parseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
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

func parseRow(row []string, lineNumber int) (Record, error) {
	var record Record

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return Record{}, fmt.Errorf("empty name at line %d", lineNumber)
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
	}
	record.Active = active

	return record, nil
}

func calculateStats(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var activeCount int
	minValue := records[0].Value
	maxValue := records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value < minValue {
			minValue = record.Value
		}
		if record.Value > maxValue {
			maxValue = record.Value
		}
		if record.Active {
			activeCount++
		}
	}

	average := sum / float64(len(records))
	return average, maxValue - minValue, activeCount
}

func filterRecords(records []Record, predicate func(Record) bool) []Record {
	var filtered []Record
	for _, record := range records {
		if predicate(record) {
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
	records, err := parseCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully parsed %d records\n", len(records))

	average, rangeValue, activeCount := calculateStats(records)
	fmt.Printf("Average value: %.2f\n", average)
	fmt.Printf("Value range: %.2f\n", rangeValue)
	fmt.Printf("Active records: %d\n", activeCount)

	highValueRecords := filterRecords(records, func(r Record) bool {
		return r.Value > 50.0
	})
	fmt.Printf("Records with value > 50: %d\n", len(highValueRecords))

	for i, record := range records {
		if i < 3 {
			fmt.Printf("Sample record %d: ID=%d, Name=%s, Value=%.2f, Active=%v\n",
				i+1, record.ID, record.Name, record.Value, record.Active)
		}
	}
}
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID    int
    Name  string
    Value float64
}

func processCSV(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
    lineNum := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
        }

        lineNum++
        if lineNum == 1 {
            continue
        }

        if len(line) != 3 {
            return nil, fmt.Errorf("invalid column count on line %d", lineNum)
        }

        id, err := strconv.Atoi(line[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID on line %d: %w", lineNum, err)
        }

        name := line[1]
        if name == "" {
            return nil, fmt.Errorf("empty name on line %d", lineNum)
        }

        value, err := strconv.ParseFloat(line[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value on line %d: %w", lineNum, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    return records, nil
}

func calculateStats(records []Record) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    var max float64 = records[0].Value

    for _, r := range records {
        sum += r.Value
        if r.Value > max {
            max = r.Value
        }
    }

    average := sum / float64(len(records))
    return average, max
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        os.Exit(1)
    }

    records, err := processCSV(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    avg, max := calculateStats(records)
    fmt.Printf("Processed %d records\n", len(records))
    fmt.Printf("Average value: %.2f\n", avg)
    fmt.Printf("Maximum value: %.2f\n", max)
}package main

import (
	"fmt"
)

func calculateMovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 || windowSize > len(data) {
		return nil
	}

	result := make([]float64, 0, len(data)-windowSize+1)

	for i := 0; i <= len(data)-windowSize; i++ {
		sum := 0.0
		for j := 0; j < windowSize; j++ {
			sum += data[i+j]
		}
		average := sum / float64(windowSize)
		result = append(result, average)
	}

	return result
}

func main() {
	sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3

	averages := calculateMovingAverage(sampleData, window)
	fmt.Printf("Moving averages with window size %d: %v\n", window, averages)
}