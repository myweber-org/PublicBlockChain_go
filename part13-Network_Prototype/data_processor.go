package main

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

// ParseUserData attempts to parse JSON into a User struct.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ParseUserData(rawData []byte) (*User, error) {
	valid, err := ValidateJSON(rawData)
	if !valid {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(rawData, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}
	return &user, nil
}

func main() {
	jsonData := []byte(`{"id": 1, "name": "Alice", "email": "alice@example.com"}`)

	user, err := ParseUserData(jsonData)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Parsed User: ID=%d, Name=%s, Email=%s\n", user.ID, user.Name, user.Email)
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataProcessor struct {
	InputPath  string
	OutputPath string
}

func NewDataProcessor(input, output string) *DataProcessor {
	return &DataProcessor{
		InputPath:  input,
		OutputPath: output,
	}
}

func (dp *DataProcessor) Process() error {
	inputFile, err := os.Open(dp.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dp.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	cleanedCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		recordCount++
		cleanedRecord := dp.cleanRecord(record)

		if dp.isValidRecord(cleanedRecord) {
			if err := writer.Write(cleanedRecord); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
			cleanedCount++
		}
	}

	fmt.Printf("Processed %d records, cleaned %d records\n", recordCount, cleanedCount)
	return nil
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
	cleaned := make([]string, len(record))
	for i, field := range record {
		cleaned[i] = strings.TrimSpace(field)
	}
	return cleaned
}

func (dp *DataProcessor) isValidRecord(record []string) bool {
	for _, field := range record {
		if field == "" {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	processor := NewDataProcessor(os.Args[1], os.Args[2])
	if err := processor.Process(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
package main

import (
	"fmt"
)

// CalculateMovingAverage computes the moving average of a slice of float64 values
// over a specified window size. Returns a slice of averages.
func CalculateMovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 || windowSize > len(data) {
		return nil
	}

	var result []float64
	var sum float64

	// Calculate initial sum for the first window
	for i := 0; i < windowSize; i++ {
		sum += data[i]
	}
	result = append(result, sum/float64(windowSize))

	// Slide the window and update the sum
	for i := windowSize; i < len(data); i++ {
		sum = sum - data[i-windowSize] + data[i]
		result = append(result, sum/float64(windowSize))
	}

	return result
}

func main() {
	// Example usage
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3

	averages := CalculateMovingAverage(data, window)
	fmt.Printf("Moving averages with window size %d: %v\n", window, averages)
}