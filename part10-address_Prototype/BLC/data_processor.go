
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

func ParseCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

    for lineNumber := 1; ; lineNumber++ {
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

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
        }

        records = append(records, DataRecord{
            ID:    id,
            Name:  row[1],
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) error {
    idSet := make(map[int]bool)
    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("invalid record ID: %d must be positive", record.ID)
        }
        if record.Name == "" {
            return fmt.Errorf("record %d has empty name", record.ID)
        }
        if record.Value < 0 {
            return fmt.Errorf("record %d has negative value: %f", record.ID, record.Value)
        }
        if idSet[record.ID] {
            return fmt.Errorf("duplicate ID found: %d", record.ID)
        }
        idSet[record.ID] = true
    }
    return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
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
package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateAndNormalize(data UserData) (UserData, error) {
	var normalized UserData

	normalized.Username = strings.TrimSpace(data.Username)
	if normalized.Username == "" {
		return UserData{}, fmt.Errorf("username cannot be empty")
	}

	normalized.Email = strings.ToLower(strings.TrimSpace(data.Email))
	if !strings.Contains(normalized.Email, "@") {
		return UserData{}, fmt.Errorf("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return UserData{}, fmt.Errorf("age must be between 0 and 150")
	}
	normalized.Age = data.Age

	return normalized, nil
}

func main() {
	testData := UserData{
		Username: "  JohnDoe  ",
		Email:    "  TEST@EXAMPLE.COM",
		Age:      25,
	}

	result, err := ValidateAndNormalize(testData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Original: %+v\n", testData)
	fmt.Printf("Normalized: %+v\n", result)
}