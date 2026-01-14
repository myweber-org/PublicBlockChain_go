
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataProcessor struct {
	FilePath string
	Headers  []string
	Records  [][]string
}

func NewDataProcessor(filePath string) *DataProcessor {
	return &DataProcessor{
		FilePath: filePath,
	}
}

func (dp *DataProcessor) LoadAndValidate() error {
	file, err := os.Open(dp.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}
	dp.Headers = headers

	dp.Records = [][]string{}
	lineNumber := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(record) != len(headers) {
			return fmt.Errorf("column count mismatch at line %d: expected %d, got %d", lineNumber, len(headers), len(record))
		}

		for i, value := range record {
			record[i] = strings.TrimSpace(value)
			if record[i] == "" {
				return fmt.Errorf("empty value at line %d, column %d", lineNumber, i+1)
			}
		}

		dp.Records = append(dp.Records, record)
		lineNumber++
	}

	if len(dp.Records) == 0 {
		return fmt.Errorf("no data records found in file")
	}

	return nil
}

func (dp *DataProcessor) GetColumnStats(columnIndex int) (min, max string, count int) {
	if columnIndex < 0 || columnIndex >= len(dp.Headers) {
		return "", "", 0
	}

	if len(dp.Records) == 0 {
		return "", "", 0
	}

	min = dp.Records[0][columnIndex]
	max = dp.Records[0][columnIndex]
	count = len(dp.Records)

	for _, record := range dp.Records {
		value := record[columnIndex]
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	return min, max, count
}

func (dp *DataProcessor) FilterRecords(predicate func([]string) bool) [][]string {
	var filtered [][]string
	for _, record := range dp.Records {
		if predicate(record) {
			filtered = append(filtered, record)
		}
	}
	return filtered
}