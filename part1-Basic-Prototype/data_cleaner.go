
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) AddItem(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return false
	}
	dc.seen[normalized] = true
	return true
}

func (dc *DataCleaner) UniqueCount() int {
	return len(dc.seen)
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	samples := []string{"  Apple  ", "apple", "BANANA", "banana ", "Cherry"}
	
	fmt.Println("Processing items:")
	for _, item := range samples {
		normalized := cleaner.Normalize(item)
		isDup := cleaner.IsDuplicate(item)
		fmt.Printf("Original: '%s' -> Normalized: '%s' -> Duplicate: %v\n", 
			item, normalized, isDup)
	}
	
	fmt.Printf("\nTotal unique items: %d\n", cleaner.UniqueCount())
	
	cleaner.Reset()
	fmt.Printf("After reset, unique items: %d\n", cleaner.UniqueCount())
}package datautils

import "sort"

func RemoveDuplicates[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return slice
	}

	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func RemoveDuplicatesSorted[T comparable](slice []T) []T {
	if len(slice) < 2 {
		return slice
	}

	sort.Slice(slice, func(i, j int) bool {
		switch v := any(slice[i]).(type) {
		case int:
			return v < any(slice[j]).(int)
		case string:
			return v < any(slice[j]).(string)
		default:
			return false
		}
	})

	j := 0
	for i := 1; i < len(slice); i++ {
		if slice[j] != slice[i] {
			j++
			slice[j] = slice[i]
		}
	}

	return slice[:j+1]
}
package main

import (
	"encoding/csv"
	"io"
	"strings"
)

func CleanCSVData(input io.Reader, output io.Writer) error {
	reader := csv.NewReader(input)
	writer := csv.NewWriter(output)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		cleanedRecord := make([]string, 0, len(record))
		hasData := false

		for _, field := range record {
			trimmed := strings.TrimSpace(field)
			cleanedRecord = append(cleanedRecord, trimmed)
			if trimmed != "" {
				hasData = true
			}
		}

		if hasData {
			if err := writer.Write(cleanedRecord); err != nil {
				return err
			}
		}
	}

	return nil
}