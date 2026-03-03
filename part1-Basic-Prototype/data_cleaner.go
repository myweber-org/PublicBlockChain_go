
package main

import (
	"regexp"
	"strings"
)

func SanitizeCSVField(input string) string {
	if input == "" {
		return input
	}

	// Remove leading/trailing whitespace
	trimmed := strings.TrimSpace(input)

	// Remove any double quotes that could break CSV formatting
	trimmed = strings.ReplaceAll(trimmed, "\"", "'")

	// Remove newlines and carriage returns
	re := regexp.MustCompile(`[\r\n]+`)
	trimmed = re.ReplaceAllString(trimmed, " ")

	// Escape existing commas if not already quoted
	if strings.Contains(trimmed, ",") && !strings.HasPrefix(trimmed, "\"") {
		trimmed = "\"" + trimmed + "\""
	}

	return trimmed
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}