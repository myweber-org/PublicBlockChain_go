
package main

import (
	"fmt"
	"strings"
)

func CleanStringSlice(input []string) []string {
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
	data := []string{"  apple", "banana  ", "apple", "", "  cherry  ", "banana"}
	cleaned := CleanStringSlice(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}