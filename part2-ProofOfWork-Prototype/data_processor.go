
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) NormalizeEmail(email string) (string, bool) {
	cleaned := dp.CleanString(email)
	normalized := strings.ToLower(cleaned)
	return normalized, dp.ValidateEmail(normalized)
}

func (dp *DataProcessor) ProcessInput(input string) (string, bool) {
	cleaned := dp.CleanString(input)
	if cleaned == "" {
		return "", false
	}
	return cleaned, true
}