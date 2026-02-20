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