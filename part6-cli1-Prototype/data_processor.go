package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateAndParseJSON(rawData []byte) (*UserData, error) {
	var user UserData
	if err := json.Unmarshal(rawData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if user.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.Name == "" {
		return nil, fmt.Errorf("user name cannot be empty")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user email cannot be empty")
	}

	return &user, nil
}

func main() {
	jsonStr := `{"id": 123, "name": "John Doe", "email": "john@example.com"}`
	user, err := ValidateAndParseJSON([]byte(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed user: %+v\n", user)
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

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)

	headerSkipped := false
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		if !headerSkipped {
			headerSkipped = true
			continue
		}

		if len(row) < 4 {
			return nil, errors.New("invalid row format")
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(row[1])

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			continue
		}

		valid := strings.ToLower(strings.TrimSpace(row[3])) == "true"

		record := DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
			Valid: valid,
		}
		records = append(records, record)
	}

	return records, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Valid && record.Value > 0 {
			filtered = append(filtered, record)
		}
	}
	return filtered
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

func ProcessDataFile(filename string) (float64, error) {
	records, err := ParseCSVFile(filename)
	if err != nil {
		return 0, err
	}

	validRecords := FilterValidRecords(records)
	if len(validRecords) == 0 {
		return 0, errors.New("no valid records found")
	}

	average := CalculateAverage(validRecords)
	return average, nil
}