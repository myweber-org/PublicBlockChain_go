
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct{}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range items {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (dc *DataCleaner) NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func main() {
	cleaner := &DataCleaner{}
	
	data := []string{"Apple", "apple", " Banana ", "banana", "Apple"}
	fmt.Println("Original:", data)
	
	cleaned := cleaner.RemoveDuplicates(data)
	fmt.Println("After deduplication:", cleaned)
	
	normalized := make([]string, len(cleaned))
	for i, item := range cleaned {
		normalized[i] = cleaner.NormalizeString(item)
	}
	fmt.Println("After normalization:", normalized)
}