package main

import (
	"strings"
)

// CleanString removes duplicate spaces and trims leading/trailing whitespace
func CleanString(input string) string {
	// Trim spaces from start and end
	trimmed := strings.TrimSpace(input)
	
	// Split by spaces and filter out empty strings
	words := strings.Fields(trimmed)
	
	// Join back with single spaces
	return strings.Join(words, " ")
}

// RemoveDuplicates removes duplicate entries from a slice of strings
func RemoveDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// CleanSlice applies CleanString to each element and removes duplicates
func CleanSlice(items []string) []string {
	cleaned := make([]string, len(items))
	
	for i, item := range items {
		cleaned[i] = CleanString(item)
	}
	
	return RemoveDuplicates(cleaned)
}