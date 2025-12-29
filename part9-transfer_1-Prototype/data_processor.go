
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (dp *DataProcessor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	cleaned := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return cleaned
}

func (dp *DataProcessor) NormalizeCase(input string, lower bool) string {
	if lower {
		return strings.ToLower(input)
	}
	return strings.ToUpper(input)
}

func (dp *DataProcessor) RemoveSpecialChars(input string, keepPattern string) string {
	if keepPattern == "" {
		keepPattern = `[^a-zA-Z0-9\s]`
	}
	re := regexp.MustCompile(keepPattern)
	return re.ReplaceAllString(input, "")
}

func main() {
	processor := NewDataProcessor()
	sample := "  Hello   World!  This  is  a  TEST.  "
	
	cleaned := processor.CleanInput(sample)
	normalized := processor.NormalizeCase(cleaned, true)
	final := processor.RemoveSpecialChars(normalized, `[^a-zA-Z\s]`)
	
	println("Original:", sample)
	println("Processed:", final)
}