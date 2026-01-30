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

func ProcessCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	_, err = reader.Read()
	if err != nil {
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

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func CalculateTotal(records []Record) float64 {
	var total float64
	for _, r := range records {
		total += r.Value
	}
	return total
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Total value: %.2f\n", CalculateTotal(records))
}package main

import (
	"errors"
	"strings"
)

type DataRecord struct {
	ID    string
	Value string
	Valid bool
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if len(record.Value) > 100 {
		return errors.New("value exceeds maximum length")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Value = strings.ToUpper(strings.TrimSpace(record.Value))
	transformed.Valid = true
	return transformed
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
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	Name   string
	Age    int
	Active bool
}

func processCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) != 3 {
			continue
		}

		age, err := strconv.Atoi(row[1])
		if err != nil {
			continue
		}

		active, err := strconv.ParseBool(row[2])
		if err != nil {
			continue
		}

		records = append(records, Record{
			Name:   row[0],
			Age:    age,
			Active: active,
		})
	}

	return records, nil
}

func calculateAverageAge(records []Record) float64 {
	if len(records) == 0 {
		return 0
	}

	total := 0
	for _, r := range records {
		total += r.Age
	}

	return float64(total) / float64(len(records))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		return
	}

	records, err := processCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Average age: %.2f\n", calculateAverageAge(records))

	activeCount := 0
	for _, r := range records {
		if r.Active {
			activeCount++
		}
	}
	fmt.Printf("Active users: %d\n", activeCount)
}