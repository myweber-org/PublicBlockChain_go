
package main

import (
    "errors"
    "strings"
)

type UserData struct {
    Username string
    Email    string
}

func ValidateUserData(data UserData) error {
    if strings.TrimSpace(data.Username) == "" {
        return errors.New("username cannot be empty")
    }
    if !strings.Contains(data.Email, "@") {
        return errors.New("invalid email format")
    }
    return nil
}

func TransformUsername(data UserData) UserData {
    data.Username = strings.ToLower(strings.TrimSpace(data.Username))
    return data
}

func ProcessUserInput(rawUsername, rawEmail string) (UserData, error) {
    userData := UserData{
        Username: rawUsername,
        Email:    rawEmail,
    }

    if err := ValidateUserData(userData); err != nil {
        return UserData{}, err
    }

    userData = TransformUsername(userData)
    return userData, nil
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
	records := []Record{}
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

		record, err := parseRow(row, lineNum)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string, lineNum int) (Record, error) {
	var record Record
	var err error

	record.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return Record{}, fmt.Errorf("empty name at line %d", lineNum)
	}

	record.Value, err = strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
	}

	record.Active, err = strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid active flag at line %d: %w", lineNum, err)
	}

	return record, nil
}

func validateRecords(records []Record) error {
	if len(records) == 0 {
		return errors.New("no records found")
	}

	idSet := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", record.ID)
		}
		if idSet[record.ID] {
			return fmt.Errorf("duplicate ID %d found", record.ID)
		}
		idSet[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("negative value %f for ID %d", record.Value, record.ID)
		}
	}

	return nil
}

func processData(filename string) error {
	records, err := parseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	if err := validateRecords(records); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	totalValue := 0.0
	activeCount := 0
	for _, record := range records {
		totalValue += record.Value
		if record.Active {
			activeCount++
		}
	}

	fmt.Printf("Processing complete:\n")
	fmt.Printf("  Total records: %d\n", len(records))
	fmt.Printf("  Active records: %d\n", activeCount)
	fmt.Printf("  Total value: %.2f\n", totalValue)
	fmt.Printf("  Average value: %.2f\n", totalValue/float64(len(records)))

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	if err := processData(os.Args[1]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}