
package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
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

import "fmt"

func removeDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{4, 2, 8, 2, 4, 9, 8, 1}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package datautils

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
package main

import (
    "fmt"
    "strings"
)

func DeduplicateStrings(slice []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    for _, item := range slice {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    return result
}

func ValidateEmail(email string) bool {
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    if len(parts[0]) == 0 || len(parts[1]) == 0 {
        return false
    }
    return true
}

func main() {
    emails := []string{
        "test@example.com",
        "user@domain.org",
        "test@example.com",
        "invalid-email",
        "another@test.net",
        "user@domain.org",
    }

    uniqueEmails := DeduplicateStrings(emails)
    fmt.Println("Unique emails:", uniqueEmails)

    for _, email := range uniqueEmails {
        if ValidateEmail(email) {
            fmt.Printf("%s is valid\n", email)
        } else {
            fmt.Printf("%s is invalid\n", email)
        }
    }
}package utils

import "strings"

func SanitizeString(input string) string {
    trimmed := strings.TrimSpace(input)
    return strings.Join(strings.Fields(trimmed), " ")
}package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
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

func (dc *DataCleaner) Deduplicate(items []string) []string {
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
	cleaned := cleaner.Deduplicate(data)
	
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
	
	cleaner.Reset()
	
	moreData := []string{"grape", "GRAPE", "kiwi"}
	moreCleaned := cleaner.Deduplicate(moreData)
	fmt.Println("Second batch:", moreCleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedCount int
	duplicateCount int
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processedCount: 0,
		duplicateCount: 0,
	}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
			dc.processedCount++
		} else {
			dc.duplicateCount++
		}
	}
	
	return result
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	
	if parts[0] == "" || parts[1] == "" {
		return false
	}
	
	return strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) GetStats() (int, int) {
	return dc.processedCount, dc.duplicateCount
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{
		"user@example.com",
		"  user@example.com  ",
		"test@domain.org",
		"invalid-email",
		"",
		"another@test.com",
		"test@domain.org",
	}
	
	uniqueEmails := cleaner.RemoveDuplicates(data)
	
	fmt.Println("Cleaned data:")
	for _, email := range uniqueEmails {
		valid := cleaner.ValidateEmail(email)
		status := "valid"
		if !valid {
			status = "invalid"
		}
		fmt.Printf("  %s (%s)\n", email, status)
	}
	
	processed, duplicates := cleaner.GetStats()
	fmt.Printf("\nStatistics:\n")
	fmt.Printf("  Processed items: %d\n", processed)
	fmt.Printf("  Duplicates found: %d\n", duplicates)
}