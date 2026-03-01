
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

func (dc *DataCleaner) CleanString(input string) string {
    trimmed := strings.TrimSpace(input)
    normalized := strings.ToLower(trimmed)
    return normalized
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
    cleaned := dc.CleanString(value)
    if dc.seen[cleaned] {
        return true
    }
    dc.seen[cleaned] = true
    return false
}

func (dc *DataCleaner) ProcessList(items []string) []string {
    var result []string
    for _, item := range items {
        if !dc.IsDuplicate(item) {
            cleaned := dc.CleanString(item)
            result = append(result, cleaned)
        }
    }
    return result
}

func main() {
    cleaner := NewDataCleaner()
    
    sampleData := []string{
        "  Apple  ",
        "apple",
        "BANANA",
        "  banana  ",
        "Orange",
        "ORANGE",
        "Grape",
    }
    
    cleaned := cleaner.ProcessList(sampleData)
    
    fmt.Println("Original count:", len(sampleData))
    fmt.Println("Cleaned count:", len(cleaned))
    fmt.Println("Unique items:", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

func CleanData(input []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	data := []string{"  apple ", "banana", "  apple", "banana ", "  ", "cherry"}
	cleaned := CleanData(data)
	fmt.Println("Cleaned data:", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

func CleanData(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	data := []string{" apple ", "banana", " apple ", "", "banana", " cherry "}
	cleaned := CleanData(data)
	fmt.Println("Cleaned data:", cleaned)
}