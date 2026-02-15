
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

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(normalized) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && !strings.ContainsAny(item, "!@#$%")
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"apple", "Apple", "banana", "", "cherry!", "banana", "date"}
	cleaned := cleaner.RemoveDuplicates(data)
	
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
	
	cleaner.Reset()
	testData := []string{"test", "test", "TEST"}
	fmt.Println("Reset test:", cleaner.RemoveDuplicates(testData))
}