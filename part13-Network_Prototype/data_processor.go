
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
        }

        if len(row) != 3 {
            continue
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            continue
        }

        name := row[1]

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            continue
        }

        records = append(records, DataRecord{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    return records, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    for _, record := range records {
        sum += record.Value
    }

    average := sum / float64(len(records))

    var variance float64
    for _, record := range records {
        diff := record.Value - average
        variance += diff * diff
    }
    variance = variance / float64(len(records))

    return average, variance
}

func FilterRecords(records []DataRecord, minValue float64) []DataRecord {
    filtered := make([]DataRecord, 0)
    for _, record := range records {
        if record.Value >= minValue {
            filtered = append(filtered, record)
        }
    }
    return filtered
}package main

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

func ProcessCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record
	line := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", line, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d", line)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", line, err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", line, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  row[1],
			Value: value,
		})
		line++
	}

	return records, nil
}

func ValidateRecords(records []Record) error {
	seen := make(map[int]bool)
	for _, r := range records {
		if r.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", r.ID)
		}
		if seen[r.ID] {
			return fmt.Errorf("duplicate ID %d found", r.ID)
		}
		if r.Value < 0 {
			return fmt.Errorf("negative value %f for ID %d", r.Value, r.ID)
		}
		seen[r.ID] = true
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}

	if err := ValidateRecords(records); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d records\n", len(records))
	for _, r := range records {
		fmt.Printf("ID: %d, Name: %s, Value: %.2f\n", r.ID, r.Name, r.Value)
	}
}