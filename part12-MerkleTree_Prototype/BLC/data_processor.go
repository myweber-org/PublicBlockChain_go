
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
	Value float64
	Valid bool
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
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
			continue
		}

		record, err := parseRecord(row)
		if err != nil {
			fmt.Printf("Warning: line %d - %v\n", lineNumber, err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRecord(row []string) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID format: %s", row[0])
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value format: %s", row[2])
	}
	record.Value = value

	validStr := strings.ToLower(strings.TrimSpace(row[3]))
	record.Valid = validStr == "true" || validStr == "1" || validStr == "yes"

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := ProcessCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total records processed: %d\n", len(records))

	validRecords := FilterValidRecords(records)
	fmt.Printf("Valid records: %d\n", len(validRecords))

	average := CalculateAverage(validRecords)
	fmt.Printf("Average value of valid records: %.2f\n", average)

	for i, record := range validRecords {
		if i < 5 {
			fmt.Printf("Record %d: ID=%d, Name=%s, Value=%.2f\n",
				i+1, record.ID, record.Name, record.Value)
		}
	}
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
    Value string
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

        if len(row) < 3 {
            continue
        }

        record := DataRecord{
            ID:    strings.TrimSpace(row[0]),
            Name:  strings.TrimSpace(row[1]),
            Value: strings.TrimSpace(row[2]),
        }

        if record.ID == "" || record.Name == "" {
            continue
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, fmt.Errorf("no valid records found in file")
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) []DataRecord {
    var validRecords []DataRecord
    seenIDs := make(map[string]bool)

    for _, record := range records {
        if seenIDs[record.ID] {
            continue
        }
        if len(record.Value) > 0 && record.Value != "null" {
            validRecords = append(validRecords, record)
            seenIDs[record.ID] = true
        }
    }

    return validRecords
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file_path>")
        return
    }

    records, err := ProcessCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        return
    }

    validRecords := ValidateRecords(records)
    fmt.Printf("Processed %d valid records from %d total records\n", 
        len(validRecords), len(records))
}