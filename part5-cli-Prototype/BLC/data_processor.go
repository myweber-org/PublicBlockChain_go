
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
	Category  string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Category == "" {
		return errors.New("category cannot be empty")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Category = strings.ToUpper(record.Category)
	transformed.Value = record.Value * 1.1
	return transformed
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

func main() {
	records := []DataRecord{
		{ID: "A001", Value: 100.0, Timestamp: time.Now(), Category: "sales"},
		{ID: "A002", Value: 250.5, Timestamp: time.Now(), Category: "inventory"},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	for _, record := range processed {
		fmt.Printf("Processed: ID=%s, Value=%.2f, Category=%s\n",
			record.ID, record.Value, record.Category)
	}
}
package main

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
	user := UserData{
		Username: rawUsername,
		Email:    rawEmail,
		Age:      rawAge,
	}

	user = TransformUsername(user)

	if err := ValidateUserData(user); err != nil {
		return UserData{}, err
	}

	return user, nil
}
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

func ValidateJSONStructure(rawData string) (bool, error) {
    var data map[string]interface{}
    decoder := json.NewDecoder(strings.NewReader(rawData))
    decoder.DisallowUnknownFields()

    if err := decoder.Decode(&data); err != nil {
        return false, fmt.Errorf("invalid JSON structure: %w", err)
    }

    if len(data) == 0 {
        return false, fmt.Errorf("JSON data is empty")
    }

    return true, nil
}

func ExtractJSONKeys(rawData string) ([]string, error) {
    var data map[string]interface{}
    if err := json.Unmarshal([]byte(rawData), &data); err != nil {
        return nil, err
    }

    keys := make([]string, 0, len(data))
    for key := range data {
        keys = append(keys, key)
    }
    return keys, nil
}

func main() {
    sampleJSON := `{"name": "test", "value": 42, "active": true}`
    
    valid, err := ValidateJSONStructure(sampleJSON)
    if err != nil {
        fmt.Printf("Validation error: %v\n", err)
        return
    }
    fmt.Printf("JSON is valid: %v\n", valid)

    keys, err := ExtractJSONKeys(sampleJSON)
    if err != nil {
        fmt.Printf("Key extraction error: %v\n", err)
        return
    }
    fmt.Printf("Extracted keys: %v\n", keys)
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
	ID    int
	Name  string
	Value float64
}

func ProcessCSVFile(filename string) ([]Record, error) {
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

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid records found in file")
	}

	return records, nil
}

func CalculateStatistics(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var max float64
	count := len(records)

	for i, record := range records {
		sum += record.Value
		if i == 0 || record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(count)
	return average, max, count
}

func ValidateRecords(records []Record) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid record ID: %d must be positive", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("negative value not allowed for record ID %d", record.ID)
		}
	}

	return nil
}