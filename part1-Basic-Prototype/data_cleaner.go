
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
}