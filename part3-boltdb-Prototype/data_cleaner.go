
package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "cherry", "banana", "date"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
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

func (dc *DataCleaner) IsDuplicate(raw string) bool {
	normalized := dc.Normalize(raw)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) CleanList(items []string) []string {
	dc.seen = make(map[string]bool)
	var cleaned []string
	for _, item := range items {
		if !dc.IsDuplicate(item) {
			cleaned = append(cleaned, dc.Normalize(item))
		}
	}
	return cleaned
}

func main() {
	cleaner := NewDataCleaner()
	data := []string{"  Apple ", "banana", "APPLE", " Banana ", "Cherry"}
	result := cleaner.CleanList(data)
	fmt.Println("Cleaned data:", result)
}