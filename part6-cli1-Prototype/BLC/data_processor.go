
package main

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

        if len(row) < 3 {
            continue
        }

        record := DataRecord{
            ID:    strings.TrimSpace(row[0]),
            Name:  strings.TrimSpace(row[1]),
            Email: strings.TrimSpace(row[2]),
            Valid: validateEmail(strings.TrimSpace(row[2])),
        }

        if record.ID != "" && record.Name != "" {
            records = append(records, record)
        }
    }

    return records, nil
}

func validateEmail(email string) bool {
    if email == "" {
        return false
    }
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func GenerateReport(records []DataRecord) {
    validCount := 0
    for _, record := range records {
        if record.Valid {
            validCount++
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Valid email addresses: %d\n", validCount)
    fmt.Printf("Invalid email addresses: %d\n", len(records)-validCount)

    if len(records) > 0 {
        fmt.Println("\nFirst 5 records:")
        for i := 0; i < len(records) && i < 5; i++ {
            status := "INVALID"
            if records[i].Valid {
                status = "VALID"
            }
            fmt.Printf("%s: %s <%s> [%s]\n", 
                records[i].ID, 
                records[i].Name, 
                records[i].Email, 
                status)
        }
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run data_processor.go <csv_file>")
        os.Exit(1)
    }

    filename := os.Args[1]
    records, err := ProcessCSVFile(filename)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    GenerateReport(records)
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

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
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

		record, err := parseRecord(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRecord(row []string, lineNumber int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return record, fmt.Errorf("empty name at line %d", lineNumber)
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}
	record.Value = value

	validStr := strings.ToLower(strings.TrimSpace(row[3]))
	if validStr != "true" && validStr != "false" {
		return record, fmt.Errorf("invalid boolean at line %d: must be 'true' or 'false'", lineNumber)
	}
	record.Valid = validStr == "true"

	return record, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Valid {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func CalculateAverage(records []DataRecord) (float64, error) {
	if len(records) == 0 {
		return 0, errors.New("no records to calculate average")
	}

	var sum float64
	count := 0
	for _, record := range records {
		if record.Valid {
			sum += record.Value
			count++
		}
	}

	if count == 0 {
		return 0, errors.New("no valid records to calculate average")
	}

	return sum / float64(count), nil
}

func FindMaxValue(records []DataRecord) (DataRecord, error) {
	if len(records) == 0 {
		return DataRecord{}, errors.New("no records to find maximum")
	}

	var maxRecord DataRecord
	found := false

	for _, record := range records {
		if record.Valid && (!found || record.Value > maxRecord.Value) {
			maxRecord = record
			found = true
		}
	}

	if !found {
		return DataRecord{}, errors.New("no valid records to find maximum")
	}

	return maxRecord, nil
}