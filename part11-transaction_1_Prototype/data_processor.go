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
    var max float64
    activeCount := 0

    for i, record := range records {
        if i == 0 || record.Value > max {
            max = record.Value
        }
        sum += record.Value
        if record.Active {
            activeCount++
        }
    }

    average := sum / float64(len(records))
    return average, max, activeCount
}

func filterRecords(records []Record, minValue float64) []Record {
    var filtered []Record
    for _, record := range records {
        if record.Value >= minValue {
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

    avg, max, activeCount := calculateStats(records)
    fmt.Printf("Total records: %d\n", len(records))
    fmt.Printf("Average value: %.2f\n", avg)
    fmt.Printf("Maximum value: %.2f\n", max)
    fmt.Printf("Active records: %d\n", activeCount)

    filtered := filterRecords(records, 50.0)
    fmt.Printf("Records with value >= 50: %d\n", len(filtered))
}package main

import (
	"errors"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(data UserData) UserData {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
	return data
}

func ProcessUserInput(rawUsername string, rawEmail string, rawAge int) (UserData, error) {
	userData := UserData{
		Username: rawUsername,
		Email:    rawEmail,
		Age:      rawAge,
	}

	userData = TransformUsername(userData)

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	return userData, nil
}