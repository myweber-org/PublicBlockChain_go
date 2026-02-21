
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
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
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

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNum, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNum)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
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

func CalculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, r := range records {
		sum += r.Value
		if r.Value > max {
			max = r.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
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

	avg, max := CalculateStats(records)
	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)

	for i, r := range records {
		if i < 3 {
			fmt.Printf("Sample record %d: ID=%d, Name=%s, Value=%.2f\n", i+1, r.ID, r.Name, r.Value)
		}
	}
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) (bool, error) {
	var js interface{}
	err := json.Unmarshal(data, &js)
	if err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}
	return true, nil
}

// ParseUserData attempts to parse JSON data into a generic map structure.
func ParseUserData(jsonData []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func main() {
	sampleJSON := []byte(`{"name": "Alice", "age": 30, "active": true}`)

	valid, err := ValidateJSON(sampleJSON)
	if err != nil {
		log.Printf("Validation error: %v", err)
	} else {
		fmt.Println("JSON is valid:", valid)
	}

	parsedData, err := ParseUserData(sampleJSON)
	if err != nil {
		log.Printf("Parsing error: %v", err)
	} else {
		fmt.Printf("Parsed data: %+v\n", parsedData)
	}
}package main

import (
	"encoding/json"
	"fmt"
)

type DataValidator struct {
	MaxSize int
}

func NewDataValidator(maxSize int) *DataValidator {
	return &DataValidator{MaxSize: maxSize}
}

func (dv *DataValidator) ValidateJSON(input []byte) (map[string]interface{}, error) {
	if len(input) > dv.MaxSize {
		return nil, fmt.Errorf("input size %d exceeds maximum allowed size %d", len(input), dv.MaxSize)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty JSON object")
	}

	return data, nil
}

func ProcessJSONData(jsonStr string) error {
	validator := NewDataValidator(1024 * 1024)
	data, err := validator.ValidateJSON([]byte(jsonStr))
	if err != nil {
		return err
	}

	fmt.Printf("Successfully validated JSON with %d fields\n", len(data))
	for key, value := range data {
		fmt.Printf("Key: %s, Type: %T\n", key, value)
	}

	return nil
}

func main() {
	sampleJSON := `{"name": "test", "value": 42, "active": true}`
	if err := ProcessJSONData(sampleJSON); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
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

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	userData = TransformUsername(userData)
	return userData, nil
}
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	cleaned := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return cleaned
}

func (dp *DataProcessor) NormalizeCase(input string) string {
	return strings.ToLower(dp.CleanString(input))
}

func (dp *DataProcessor) RemoveSpecialChars(input string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	return reg.ReplaceAllString(input, "")
}

func (dp *DataProcessor) Process(input string) string {
	cleaned := dp.CleanString(input)
	normalized := dp.NormalizeCase(cleaned)
	return dp.RemoveSpecialChars(normalized)
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
			return errors.New("username can only contain letters, digits, and underscores")
		}
	}

	if !strings.Contains(data.Email, "@") || !strings.Contains(data.Email, ".") {
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

	transformed := UserData{
		Username: NormalizeUsername(data.Username),
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Age:      data.Age,
	}

	return transformed, nil
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

func processCSV(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
    lineNum := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
        }

        lineNum++
        if lineNum == 1 {
            continue
        }

        if len(line) != 3 {
            return nil, fmt.Errorf("invalid column count on line %d", lineNum)
        }

        id, err := strconv.Atoi(line[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID on line %d: %w", lineNum, err)
        }

        name := line[1]
        if name == "" {
            return nil, fmt.Errorf("empty name on line %d", lineNum)
        }

        value, err := strconv.ParseFloat(line[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value on line %d: %w", lineNum, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    return records, nil
}

func calculateStats(records []Record) (float64, float64) {
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

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        os.Exit(1)
    }

    records, err := processCSV(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    avg, max := calculateStats(records)
    fmt.Printf("Processed %d records\n", len(records))
    fmt.Printf("Average value: %.2f\n", avg)
    fmt.Printf("Maximum value: %.2f\n", max)
}