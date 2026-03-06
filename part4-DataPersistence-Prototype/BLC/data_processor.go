package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitized := re.ReplaceAllString(input, "")
	return strings.TrimSpace(sanitized)
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TrimAndLower(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}package main

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ParseCSVData(reader io.Reader) ([]DataRecord, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataRecord
	for i, row := range records {
		if len(row) < 4 {
			continue
		}

		record, err := validateRow(row)
		if err != nil {
			continue
		}

		data = append(data, record)
	}

	return data, nil
}

func validateRow(row []string) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, errors.New("invalid id")
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return record, errors.New("empty name")
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, errors.New("invalid value")
	}
	record.Value = value

	valid, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		record.Valid = false
	} else {
		record.Valid = valid
	}

	return record, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if record.Valid {
			valid = append(valid, record)
		}
	}
	return valid
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}

	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	return sum / float64(len(records))
}package data

import (
	"regexp"
	"strings"
)

// ValidateEmail checks if the provided string is a valid email address.
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// SanitizeInput removes leading and trailing whitespace from a string.
func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

// ConvertToUpper transforms a string to uppercase.
func ConvertToUpper(s string) string {
	return strings.ToUpper(s)
}

// IsNumeric checks if a string contains only numeric characters.
func IsNumeric(s string) bool {
	pattern := `^[0-9]+$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}
package main

import "fmt"

func calculateAverage(numbers []int) float64 {
    if len(numbers) == 0 {
        return 0
    }
    
    sum := 0
    for _, num := range numbers {
        sum += num
    }
    
    return float64(sum) / float64(len(numbers))
}

func main() {
    data := []int{10, 20, 30, 40, 50}
    avg := calculateAverage(data)
    fmt.Printf("Average: %.2f\n", avg)
}
package main

import (
	"fmt"
)

// CalculateMovingAverage returns a slice containing the moving average of the input slice.
// The windowSize parameter defines the number of elements to average over.
// If windowSize is greater than the length of data, an empty slice is returned.
func CalculateMovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 || windowSize > len(data) {
		return []float64{}
	}

	result := make([]float64, len(data)-windowSize+1)
	for i := 0; i <= len(data)-windowSize; i++ {
		sum := 0.0
		for j := 0; j < windowSize; j++ {
			sum += data[i+j]
		}
		result[i] = sum / float64(windowSize)
	}
	return result
}

func main() {
	sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3
	averages := CalculateMovingAverage(sampleData, window)
	fmt.Printf("Moving averages (window=%d): %v\n", window, averages)
}
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) NormalizeEmail(email string) string {
	cleaned := dp.CleanString(email)
	return strings.ToLower(cleaned)
}

func (dp *DataProcessor) ProcessUserInput(name, email string) (string, string, bool) {
	cleanName := dp.CleanString(name)
	normalizedEmail := dp.NormalizeEmail(email)
	isValid := dp.ValidateEmail(normalizedEmail)
	
	return cleanName, normalizedEmail, isValid
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
	ID    string
	Name  string
	Email string
	Valid bool
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
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

		if len(row) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Name:  strings.TrimSpace(row[1]),
			Email: strings.TrimSpace(row[2]),
			Valid: validateRecord(row[0], row[1], row[2]),
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(id, name, email string) bool {
	if id == "" || name == "" || email == "" {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	return true
}

func GenerateReport(records []DataRecord) {
	validCount := 0
	invalidCount := 0

	fmt.Println("=== DATA PROCESSING REPORT ===")
	for _, record := range records {
		if record.Valid {
			validCount++
			fmt.Printf("✓ Valid: %s - %s\n", record.ID, record.Name)
		} else {
			invalidCount++
			fmt.Printf("✗ Invalid: %s - %s (%s)\n", record.ID, record.Name, record.Email)
		}
	}

	fmt.Printf("\nSummary: %d valid, %d invalid, %d total\n", 
		validCount, invalidCount, len(records))
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func ProcessUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !ValidateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = SanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email":"test@example.com","username":"  JohnDoe  ","age":25}`
	processedData, err := ProcessUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed Data: %+v\n", processedData)
}package main

import (
    "encoding/csv"
    "errors"
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

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
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
            return nil, err
        }

        if len(row) < 4 {
            continue
        }

        record, parseErr := parseRow(row)
        if parseErr == nil {
            records = append(records, record)
        }
    }

    return records, nil
}

func parseRow(row []string) (DataRecord, error) {
    var record DataRecord
    var err error

    record.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return record, errors.New("invalid ID format")
    }

    record.Name = strings.TrimSpace(row[1])
    if record.Name == "" {
        return record, errors.New("name cannot be empty")
    }

    record.Value, err = strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return record, errors.New("invalid value format")
    }

    validStr := strings.ToLower(strings.TrimSpace(row[3]))
    record.Valid = validStr == "true" || validStr == "1"

    return record, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
    filtered := make([]DataRecord, 0)
    for _, record := range records {
        if record.Valid {
            filtered = append(filtered, record)
        }
    }
    return filtered
}

func CalculateAverage(records []DataRecord) float64 {
    if len(records) == 0 {
        return 0
    }

    total := 0.0
    for _, record := range records {
        total += record.Value
    }
    return total / float64(len(records))
}
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) NormalizeEmail(email string) (string, bool) {
	cleaned := dp.CleanString(email)
	lowerEmail := strings.ToLower(cleaned)
	return lowerEmail, dp.ValidateEmail(lowerEmail)
}
package data_processor

import (
	"regexp"
	"strings"
)

// CleanInput removes extra whitespace and trims the given string
func CleanInput(input string) string {
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")
	return strings.TrimSpace(cleaned)
}

// NormalizeEmail converts email to lowercase and trims spaces
func NormalizeEmail(email string) string {
	return strings.ToLower(CleanInput(email))
}

// ValidateUsername checks if username contains only allowed characters
func ValidateUsername(username string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]{3,20}$`, username)
	return matched
}