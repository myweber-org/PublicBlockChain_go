
package main

import (
	"fmt"
	"math"
)

// FilterAndTransform processes a slice of integers, filters out values below threshold,
// and applies a transformation (square root of absolute value).
func FilterAndTransform(numbers []int, threshold int) []float64 {
	var result []float64
	for _, num := range numbers {
		if num > threshold {
			transformed := math.Sqrt(math.Abs(float64(num)))
			result = append(result, transformed)
		}
	}
	return result
}

func main() {
	input := []int{-10, 5, 3, 15, 8, -2, 25}
	threshold := 5
	output := FilterAndTransform(input, threshold)
	fmt.Printf("Processed slice: %v\n", output)
}
package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Comments string
}

func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(trimmed, "")
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ProcessUserData(data UserData) (UserData, error) {
	sanitizedData := UserData{
		Username: SanitizeInput(data.Username),
		Email:    SanitizeInput(data.Email),
		Comments: SanitizeInput(data.Comments),
	}

	if !ValidateEmail(sanitizedData.Email) {
		return sanitizedData, &InvalidEmailError{Email: sanitizedData.Email}
	}

	if len(sanitizedData.Username) < 3 {
		return sanitizedData, &InvalidUsernameError{Username: sanitizedData.Username}
	}

	return sanitizedData, nil
}

type InvalidEmailError struct {
	Email string
}

func (e *InvalidEmailError) Error() string {
	return "Invalid email format: " + e.Email
}

type InvalidUsernameError struct {
	Username string
}

func (e *InvalidUsernameError) Error() string {
	return "Username must be at least 3 characters long: " + e.Username
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
    lineNumber := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
        }

        lineNumber++
        if lineNumber == 1 {
            continue
        }

        if len(line) != 3 {
            return nil, errors.New("invalid column count")
        }

        id, err := strconv.Atoi(line[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID format: %w", err)
        }

        value, err := strconv.ParseFloat(line[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value format: %w", err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  line[1],
            Value: value,
        })
    }

    return records, nil
}

func ValidateRecords(records []Record) error {
    if len(records) == 0 {
        return errors.New("no records to validate")
    }

    seen := make(map[int]bool)
    for _, r := range records {
        if seen[r.ID] {
            return fmt.Errorf("duplicate ID found: %d", r.ID)
        }
        seen[r.ID] = true

        if r.Value < 0 {
            return fmt.Errorf("negative value for ID %d: %f", r.ID, r.Value)
        }
    }

    return nil
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
}package main

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

func ProcessUserInput(rawData UserData) (UserData, error) {
	if err := ValidateUserData(rawData); err != nil {
		return UserData{}, err
	}
	processedData := TransformUserName(rawData)
	return processedData, nil
}
package main

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func TransformUserData(rawEmail, rawUsername string) (string, string, error) {
	if err := ValidateEmail(rawEmail); err != nil {
		return "", "", err
	}
	normalizedEmail := strings.ToLower(strings.TrimSpace(rawEmail))
	normalizedUsername := NormalizeUsername(rawUsername)
	return normalizedEmail, normalizedUsername, nil
}package data_processor

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func ParseCSV(reader io.Reader) ([]Record, error) {
	csvReader := csv.NewReader(reader)
	records := make([]Record, 0)

	// Skip header
	if _, err := csvReader.Read(); err != nil {
		return nil, err
	}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) != 3 {
			return nil, errors.New("invalid csv format")
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, err
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, err
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
	for _, rec := range records {
		if rec.ID <= 0 {
			return errors.New("invalid id")
		}
		if rec.Name == "" {
			return errors.New("empty name")
		}
		if rec.Value < 0 {
			return errors.New("negative value")
		}
		if seen[rec.ID] {
			return errors.New("duplicate id")
		}
		seen[rec.ID] = true
	}
	return nil
}
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TransformToUpper(input string) string {
	return strings.ToUpper(input)
}

func PrettyPrintJSON(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func main() {
	email := "test@example.com"
	fmt.Printf("Email %s valid: %v\n", email, ValidateEmail(email))

	str := "hello world"
	fmt.Printf("Original: %s, Transformed: %s\n", str, TransformToUpper(str))

	sample := map[string]interface{}{
		"name":  "John",
		"age":   30,
		"email": "john@example.com",
	}
	pretty, _ := PrettyPrintJSON(sample)
	fmt.Println("Pretty JSON:")
	fmt.Println(pretty)
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

    if len(records) == 0 {
        return nil, fmt.Errorf("no valid records found in file")
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) []string {
    var errors []string
    emailSet := make(map[string]bool)

    for i, record := range records {
        if record.Active != "true" && record.Active != "false" {
            errors = append(errors, fmt.Sprintf("record %d: invalid active status '%s'", i+1, record.Active))
        }

        if emailSet[record.Email] {
            errors = append(errors, fmt.Sprintf("record %d: duplicate email '%s'", i+1, record.Email))
        }
        emailSet[record.Email] = true
    }

    return errors
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        os.Exit(1)
    }

    filename := os.Args[1]
    records, err := ProcessCSVFile(filename)
    if err != nil {
        fmt.Printf("Processing error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %d records\n", len(records))

    validationErrors := ValidateRecords(records)
    if len(validationErrors) > 0 {
        fmt.Println("Validation errors found:")
        for _, errMsg := range validationErrors {
            fmt.Println("  -", errMsg)
        }
        os.Exit(1)
    }

    fmt.Println("All records validated successfully")
}