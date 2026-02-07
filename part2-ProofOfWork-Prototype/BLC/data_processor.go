package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ProcessRecord(record DataRecord) (string, error) {
	if record.ID <= 0 {
		return "", errors.New("invalid record ID")
	}

	if strings.TrimSpace(record.Name) == "" {
		return "", errors.New("record name cannot be empty")
	}

	if record.Value < 0 {
		return "", errors.New("record value cannot be negative")
	}

	processedName := strings.ToUpper(record.Name)
	formattedValue := fmt.Sprintf("%.2f", record.Value)

	result := fmt.Sprintf("Processed: ID=%d, NAME=%s, VALUE=%s",
		record.ID, processedName, formattedValue)

	return result, nil
}

func ValidateAndProcess(records []DataRecord) ([]string, []error) {
	var results []string
	var errs []error

	for _, record := range records {
		result, err := ProcessRecord(record)
		if err != nil {
			errs = append(errs, fmt.Errorf("record %d: %w", record.ID, err))
			continue
		}
		results = append(results, result)
	}

	return results, errs
}

func main() {
	records := []DataRecord{
		{ID: 1, Name: "record_one", Value: 100.50},
		{ID: 2, Name: "record_two", Value: -5.0},
		{ID: 0, Name: "record_three", Value: 75.25},
		{ID: 4, Name: "", Value: 200.0},
		{ID: 5, Name: "record_five", Value: 300.75},
	}

	results, errs := ValidateAndProcess(records)

	fmt.Println("Processing Results:")
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\nErrors:")
	for _, err := range errs {
		fmt.Println(err)
	}
}package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Email     string
	Username  string
	Age       int
}

func ValidateUserData(data UserData) error {
	if data.Age < 13 || data.Age > 120 {
		return errors.New("age must be between 13 and 120")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}

	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(data.Username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}

	return nil
}

func NormalizeUserData(data UserData) UserData {
	return UserData{
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Username: strings.TrimSpace(data.Username),
		Age:      data.Age,
	}
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	data := UserData{
		Email:    email,
		Username: username,
		Age:      age,
	}

	normalizedData := NormalizeUserData(data)
	err := ValidateUserData(normalizedData)
	if err != nil {
		return UserData{}, err
	}

	return normalizedData, nil
}
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

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

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

        records = append(records, DataRecord{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    var max float64 = records[0].Value

    for _, record := range records {
        sum += record.Value
        if record.Value > max {
            max = record.Value
        }
    }

    average := sum / float64(len(records))
    return average, max
}