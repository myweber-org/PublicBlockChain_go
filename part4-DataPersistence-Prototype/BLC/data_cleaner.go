
package main

import (
	"strings"
)

func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func TrimAndClean(items []string) []string {
	cleaned := []string{}
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

func ProcessStringList(input []string) []string {
	trimmed := TrimAndClean(input)
	deduped := RemoveDuplicates(trimmed)
	return deduped
}