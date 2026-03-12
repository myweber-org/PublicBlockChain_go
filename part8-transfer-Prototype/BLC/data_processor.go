
package main

import (
	"regexp"
	"strings"
)

// DataProcessor handles cleaning and normalization of string data
type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
}

// NewDataProcessor creates a new DataProcessor instance
func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

// CleanString removes extra whitespace and trims the input
func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	cleaned := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return cleaned
}

// NormalizeCase converts string to lowercase with first letter capitalized
func (dp *DataProcessor) NormalizeCase(input string) string {
	cleaned := dp.CleanString(input)
	if len(cleaned) == 0 {
		return cleaned
	}
	return strings.ToUpper(cleaned[:1]) + strings.ToLower(cleaned[1:])
}

// ContainsOnlyLetters checks if string contains only alphabetical characters
func (dp *DataProcessor) ContainsOnlyLetters(input string) bool {
	cleaned := dp.CleanString(input)
	for _, r := range cleaned {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == ' ') {
			return false
		}
	}
	return true
}