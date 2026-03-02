package main

import (
	"fmt"
	"strings"
)

func CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	words := strings.Fields(trimmed)
	seen := make(map[string]bool)
	var result []string

	for _, word := range words {
		if !seen[word] {
			seen[word] = true
			result = append(result, word)
		}
	}
	return strings.Join(result, " ")
}

func main() {
	testData := "  apple   banana apple   cherry banana  "
	cleaned := CleanString(testData)
	fmt.Printf("Original: '%s'\n", testData)
	fmt.Printf("Cleaned:  '%s'\n", cleaned)
}