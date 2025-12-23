package main

import (
	"errors"
	"strings"
)

type UserData struct {
	ID    int
	Name  string
	Email string
}

func ValidateUserData(data UserData) error {
	if data.ID <= 0 {
		return errors.New("invalid user ID")
	}
	if strings.TrimSpace(data.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func TransformUserName(data UserData) UserData {
	data.Name = strings.ToUpper(strings.TrimSpace(data.Name))
	return data
}

func ProcessUserInput(rawName, rawEmail string, id int) (UserData, error) {
	user := UserData{
		ID:    id,
		Name:  rawName,
		Email: rawEmail,
	}

	if err := ValidateUserData(user); err != nil {
		return UserData{}, err
	}

	user = TransformUserName(user)
	return user, nil
}
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	allowedPattern *regexp.Regexp
}

func NewDataProcessor(allowedPattern string) (*DataProcessor, error) {
	compiled, err := regexp.Compile(allowedPattern)
	if err != nil {
		return nil, err
	}
	return &DataProcessor{allowedPattern: compiled}, nil
}

func (dp *DataProcessor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.allowedPattern.FindString(trimmed)
}

func (dp *DataProcessor) ValidateInput(input string) bool {
	return dp.allowedPattern.MatchString(strings.TrimSpace(input))
}

func (dp *DataProcessor) ProcessBatch(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		cleaned := dp.CleanInput(input)
		if cleaned != "" {
			results = append(results, cleaned)
		}
	}
	return results
}
package main

import (
	"encoding/csv"
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

		if record.ID == "" || record.Name == "" {
			return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
		}

		if !strings.Contains(record.Email, "@") {
			return nil, fmt.Errorf("invalid email format at line %d", lineNumber)
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if record.Active == "true" && len(record.Name) > 0 {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	validRecords := ValidateRecords(records)
	fmt.Printf("Processed %d records, %d valid active records\n", len(records), len(validRecords))

	for i, record := range validRecords {
		if i < 3 {
			fmt.Printf("ID: %s, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
		}
	}
}