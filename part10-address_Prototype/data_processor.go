
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Value = record.Value * 1.1
	transformed.Tags = append(record.Tags, "processed")
	transformed.Tags = normalizeTags(transformed.Tags)
	return transformed
}

func normalizeTags(tags []string) []string {
	unique := make(map[string]bool)
	var result []string
	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" && !unique[normalized] {
			unique[normalized] = true
			result = append(result, normalized)
		}
	}
	return result
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record))
	}
	return processed, nil
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
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

type DataRecord struct {
    ID      int
    Name    string
    Value   float64
    Active  bool
}

type DataProcessor struct {
    records []DataRecord
}

func NewDataProcessor() *DataProcessor {
    return &DataProcessor{
        records: make([]DataRecord, 0),
    }
}

func (dp *DataProcessor) LoadFromCSV(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    lineNumber := 0
    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if lineNumber == 1 {
            continue
        }

        record, err := parseCSVRow(row)
        if err != nil {
            return fmt.Errorf("parse error at line %d: %w", lineNumber, err)
        }

        dp.records = append(dp.records, record)
    }

    return nil
}

func parseCSVRow(row []string) (DataRecord, error) {
    if len(row) != 4 {
        return DataRecord{}, errors.New("invalid number of columns")
    }

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid ID: %w", err)
    }

    name := strings.TrimSpace(row[1])
    if name == "" {
        return DataRecord{}, errors.New("name cannot be empty")
    }

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid value: %w", err)
    }

    active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid active flag: %w", err)
    }

    return DataRecord{
        ID:     id,
        Name:   name,
        Value:  value,
        Active: active,
    }, nil
}

func (dp *DataProcessor) FilterActive() []DataRecord {
    var activeRecords []DataRecord
    for _, record := range dp.records {
        if record.Active {
            activeRecords = append(activeRecords, record)
        }
    }
    return activeRecords
}

func (dp *DataProcessor) CalculateTotal() float64 {
    var total float64
    for _, record := range dp.records {
        total += record.Value
    }
    return total
}

func (dp *DataProcessor) FindByName(name string) *DataRecord {
    for _, record := range dp.records {
        if strings.EqualFold(record.Name, name) {
            return &record
        }
    }
    return nil
}

func (dp *DataProcessor) ExportToCSV(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    header := []string{"ID", "Name", "Value", "Active"}
    if err := writer.Write(header); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    for _, record := range dp.records {
        row := []string{
            strconv.Itoa(record.ID),
            record.Name,
            strconv.FormatFloat(record.Value, 'f', 2, 64),
            strconv.FormatBool(record.Active),
        }
        if err := writer.Write(row); err != nil {
            return fmt.Errorf("failed to write row: %w", err)
        }
    }

    return nil
}

func main() {
    processor := NewDataProcessor()

    if err := processor.LoadFromCSV("input.csv"); err != nil {
        fmt.Printf("Error loading data: %v\n", err)
        return
    }

    fmt.Printf("Loaded %d records\n", len(processor.records))
    fmt.Printf("Total value: %.2f\n", processor.CalculateTotal())

    activeRecords := processor.FilterActive()
    fmt.Printf("Active records: %d\n", len(activeRecords))

    if err := processor.ExportToCSV("output.csv"); err != nil {
        fmt.Printf("Error exporting data: %v\n", err)
    } else {
        fmt.Println("Data exported successfully")
    }
}