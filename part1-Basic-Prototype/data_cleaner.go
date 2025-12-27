
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

func (dc *DataCleaner) Deduplicate(values []string) []string {
	dc.seen = make(map[string]bool)
	var result []string
	for _, v := range values {
		if !dc.IsDuplicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"Apple", " apple ", "BANANA", "banana", "Cherry", "cherry "}
	
	fmt.Println("Original data:", data)
	
	deduped := cleaner.Deduplicate(data)
	fmt.Println("Deduplicated data:", deduped)
	
	testValue := "  APPLE  "
	fmt.Printf("Is '%s' duplicate? %v\n", testValue, cleaner.IsDuplicate(testValue))
}