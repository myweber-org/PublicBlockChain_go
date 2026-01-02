
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	Data []string
}

func NewDataCleaner(data []string) *DataCleaner {
	return &DataCleaner{Data: data}
}

func (dc *DataCleaner) RemoveDuplicates() []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range dc.Data {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; !exists {
			seen[trimmed] = struct{}{}
			result = append(result, trimmed)
		}
	}
	return result
}

func (dc *DataCleaner) Clean() []string {
	return dc.RemoveDuplicates()
}

func main() {
	rawData := []string{"  apple ", "banana", "  apple", "cherry  ", "banana", "", "  "}
	cleaner := NewDataCleaner(rawData)
	cleaned := cleaner.Clean()
	fmt.Println("Cleaned data:", cleaned)
}