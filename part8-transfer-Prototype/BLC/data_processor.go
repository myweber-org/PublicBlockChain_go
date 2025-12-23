
package main

import (
	"fmt"
)

// CalculateMovingAverage computes the moving average of a slice of float64 values.
// It returns a new slice where each element is the average of the preceding 'windowSize' elements.
// If the window size is larger than the number of elements up to the current index,
// it averages the available elements.
func CalculateMovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 {
		return []float64{}
	}

	result := make([]float64, len(data))
	for i := range data {
		start := i - windowSize + 1
		if start < 0 {
			start = 0
		}
		sum := 0.0
		count := 0
		for j := start; j <= i; j++ {
			sum += data[j]
			count++
		}
		result[i] = sum / float64(count)
	}
	return result
}

func main() {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3
	averages := CalculateMovingAverage(data, window)
	fmt.Printf("Data: %v\n", data)
	fmt.Printf("Moving Averages (window=%d): %v\n", window, averages)
}
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

func ParseCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}

		if len(row) < 3 {
			return nil, fmt.Errorf("invalid row format: %v", row)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %w", err)
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value format: %w", err)
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid ID: %d", record.ID)
		}

		if record.Name == "" {
			return fmt.Errorf("empty name for ID: %d", record.ID)
		}

		if record.Value < 0 {
			return fmt.Errorf("negative value for ID: %d", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID: %d", record.ID)
		}
		seenIDs[record.ID] = true
	}

	return nil
}

func CalculateTotalValue(records []DataRecord) float64 {
	var total float64
	for _, record := range records {
		total += record.Value
	}
	return total
}