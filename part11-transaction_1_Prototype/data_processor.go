package main

import (
	"encoding/csv"
	"errors"
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

func ReadCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)

	// Skip header
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 3 {
			return nil, errors.New("invalid CSV format")
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID: %v", err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value: %v", err)
		}

		record := DataRecord{
			ID:    id,
			Name:  row[1],
			Value: value,
		}
		records = append(records, record)
	}

	return records, nil
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

func FilterByThreshold(records []DataRecord, threshold float64) []DataRecord {
	filtered := make([]DataRecord, 0)
	for _, record := range records {
		if record.Value >= threshold {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func WriteProcessedData(filename string, records []DataRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Name", "Value"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range records {
		row := []string{
			strconv.Itoa(record.ID),
			record.Name,
			strconv.FormatFloat(record.Value, 'f', 2, 64),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
package main

import (
    "encoding/csv"
    "errors"
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

func ParseCSVFile(filePath string) ([]DataRecord, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    var records []DataRecord
    lineNumber := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        lineNumber++
        if lineNumber == 1 {
            continue
        }

        record, err := parseCSVLine(line, lineNumber)
        if err != nil {
            return nil, err
        }

        records = append(records, record)
    }

    return records, nil
}

func parseCSVLine(fields []string, lineNum int) (DataRecord, error) {
    if len(fields) != 4 {
        return DataRecord{}, errors.New("invalid field count at line " + strconv.Itoa(lineNum))
    }

    id, err := strconv.Atoi(strings.TrimSpace(fields[0]))
    if err != nil {
        return DataRecord{}, errors.New("invalid ID format at line " + strconv.Itoa(lineNum))
    }

    name := strings.TrimSpace(fields[1])
    if name == "" {
        return DataRecord{}, errors.New("empty name at line " + strconv.Itoa(lineNum))
    }

    value, err := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
    if err != nil {
        return DataRecord{}, errors.New("invalid value format at line " + strconv.Itoa(lineNum))
    }

    active, err := strconv.ParseBool(strings.TrimSpace(fields[3]))
    if err != nil {
        return DataRecord{}, errors.New("invalid active flag at line " + strconv.Itoa(lineNum))
    }

    return DataRecord{
        ID:     id,
        Name:   name,
        Value:  value,
        Active: active,
    }, nil
}

func FilterActiveRecords(records []DataRecord) []DataRecord {
    var filtered []DataRecord
    for _, record := range records {
        if record.Active {
            filtered = append(filtered, record)
        }
    }
    return filtered
}

func CalculateTotalValue(records []DataRecord) float64 {
    var total float64
    for _, record := range records {
        total += record.Value
    }
    return total
}

func ValidateRecord(record DataRecord) error {
    if record.ID <= 0 {
        return errors.New("ID must be positive")
    }
    if record.Name == "" {
        return errors.New("name cannot be empty")
    }
    if record.Value < 0 {
        return errors.New("value cannot be negative")
    }
    return nil
}