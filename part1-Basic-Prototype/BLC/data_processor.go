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

func processCSVFile(filename string) ([]Record, error) {
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

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func calculateStatistics(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
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

	filename := os.Args[1]
	records, err := processCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d records\n", len(records))
	
	avg, max := calculateStatistics(records)
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)
}
package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID      string
    Name    string
    Email   string
    Active  string
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
    file, err := os.Open(filePath)
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

        if lineNumber == 1 {
            continue
        }

        if len(row) < 4 {
            return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
        }

        record := DataRecord{
            ID:     strings.TrimSpace(row[0]),
            Name:   strings.TrimSpace(row[1]),
            Email:  strings.TrimSpace(row[2]),
            Active: strings.TrimSpace(row[3]),
        }

        if err := validateRecord(record); err != nil {
            return nil, fmt.Errorf("validation failed at line %d: %w", lineNumber, err)
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func validateRecord(record DataRecord) error {
    if record.ID == "" {
        return errors.New("ID cannot be empty")
    }
    if record.Name == "" {
        return errors.New("name cannot be empty")
    }
    if !strings.Contains(record.Email, "@") {
        return errors.New("invalid email format")
    }
    if record.Active != "true" && record.Active != "false" {
        return errors.New("active field must be 'true' or 'false'")
    }
    return nil
}

func GenerateSummary(records []DataRecord) map[string]int {
    summary := make(map[string]int)
    summary["total_records"] = len(records)

    activeCount := 0
    for _, record := range records {
        if record.Active == "true" {
            activeCount++
        }
    }
    summary["active_records"] = activeCount
    summary["inactive_records"] = len(records) - activeCount

    return summary
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run data_processor.go <csv_file>")
        os.Exit(1)
    }

    filePath := os.Args[1]
    records, err := ProcessCSVFile(filePath)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    summary := GenerateSummary(records)
    fmt.Printf("Successfully processed %d records\n", summary["total_records"])
    fmt.Printf("Active records: %d\n", summary["active_records"])
    fmt.Printf("Inactive records: %d\n", summary["inactive_records"])

    fmt.Println("\nFirst 3 records:")
    for i := 0; i < len(records) && i < 3; i++ {
        fmt.Printf("%+v\n", records[i])
    }
}