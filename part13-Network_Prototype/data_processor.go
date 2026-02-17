
package main

import (
    "encoding/csv"
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

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
        }

        if len(row) != 3 {
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

        records = append(records, DataRecord{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    return records, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    for _, record := range records {
        sum += record.Value
    }

    average := sum / float64(len(records))

    var variance float64
    for _, record := range records {
        diff := record.Value - average
        variance += diff * diff
    }
    variance = variance / float64(len(records))

    return average, variance
}

func FilterRecords(records []DataRecord, minValue float64) []DataRecord {
    filtered := make([]DataRecord, 0)
    for _, record := range records {
        if record.Value >= minValue {
            filtered = append(filtered, record)
        }
    }
    return filtered
}package main

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
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record
	line := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", line, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d", line)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", line, err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", line, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  row[1],
			Value: value,
		})
		line++
	}

	return records, nil
}

func ValidateRecords(records []Record) error {
	seen := make(map[int]bool)
	for _, r := range records {
		if r.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", r.ID)
		}
		if seen[r.ID] {
			return fmt.Errorf("duplicate ID %d found", r.ID)
		}
		if r.Value < 0 {
			return fmt.Errorf("negative value %f for ID %d", r.Value, r.ID)
		}
		seen[r.ID] = true
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}

	if err := ValidateRecords(records); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d records\n", len(records))
	for _, r := range records {
		fmt.Printf("ID: %d, Name: %s, Value: %.2f\n", r.ID, r.Name, r.Value)
	}
}package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func LoadCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record
	line := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		line++

		if len(row) != 3 {
			return nil, errors.New("invalid column count at line " + strconv.Itoa(line))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, errors.New("invalid ID at line " + strconv.Itoa(line))
		}

		name := row[1]
		if name == "" {
			return nil, errors.New("empty name at line " + strconv.Itoa(line))
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, errors.New("invalid value at line " + strconv.Itoa(line))
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func ValidateRecords(records []Record) error {
	seen := make(map[int]bool)
	for _, r := range records {
		if seen[r.ID] {
			return errors.New("duplicate ID found: " + strconv.Itoa(r.ID))
		}
		seen[r.ID] = true

		if r.Value < 0 {
			return errors.New("negative value for ID: " + strconv.Itoa(r.ID))
		}
	}
	return nil
}

func CalculateTotal(records []Record) float64 {
	total := 0.0
	for _, r := range records {
		total += r.Value
	}
	return total
}package main

import (
	"errors"
	"strings"
	"unicode"
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

	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	for _, r := range data.Username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return errors.New("username can only contain letters, digits and underscores")
		}
	}

	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}

	if data.Age < 18 || data.Age > 120 {
		return errors.New("age must be between 18 and 120")
	}

	return nil
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func TransformUserData(data UserData) (UserData, error) {
	if err := ValidateUserData(data); err != nil {
		return UserData{}, err
	}

	normalizedUsername := NormalizeUsername(data.Username)

	return UserData{
		Username: normalizedUsername,
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Age:      data.Age,
	}, nil
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
		if record.Active == "true" && len(record.Email) > 0 {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func GenerateReport(records []DataRecord) {
	fmt.Printf("Total records processed: %d\n", len(records))
	activeCount := 0
	for _, record := range records {
		if record.Active == "true" {
			activeCount++
		}
	}
	fmt.Printf("Active records: %d\n", activeCount)
	fmt.Printf("Inactive records: %d\n", len(records)-activeCount)
}package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatJSONString attempts to parse a JSON string and returns it in a pretty-printed format.
// If the input is invalid JSON, it returns an error message.
func FormatJSONString(input string) (string, error) {
	var data interface{}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(prettyJSON), nil
}

// ValidateJSON checks if the provided string is valid JSON.
func ValidateJSON(input string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(input), &js) == nil
}

func main() {
	// Example usage
	testString := `{"name":"test","value":42,"active":true}`
	fmt.Printf("Valid JSON: %v\n", ValidateJSON(testString))

	formatted, err := FormatJSONString(testString)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Formatted JSON:\n%s\n", formatted)
	}

	// Test with invalid JSON
	invalidString := `{"name": test}`
	fmt.Printf("\nValid JSON: %v\n", ValidateJSON(invalidString))
	_, err = FormatJSONString(invalidString)
	if err != nil {
		fmt.Printf("Expected error: %v\n", strings.TrimSpace(err.Error()))
	}
}