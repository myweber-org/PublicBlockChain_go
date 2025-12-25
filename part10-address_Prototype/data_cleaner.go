package main

import (
	"strings"
)

func TrimWhitespaceFromSlice(input []string) []string {
	output := make([]string, 0, len(input))
	for _, s := range input {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			output = append(output, trimmed)
		}
	}
	return output
}