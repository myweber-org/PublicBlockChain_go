
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

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return normalized
}

func (dp *DataProcessor) NormalizeCase(input string, toUpper bool) string {
	cleaned := dp.CleanString(input)
	if toUpper {
		return strings.ToUpper(cleaned)
	}
	return strings.ToLower(cleaned)
}

func (dp *DataProcessor) ExtractAlphanumeric(input string) string {
	alnumRegex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleaned := dp.CleanString(input)
	return alnumRegex.ReplaceAllString(cleaned, "")
}

func main() {
	processor := NewDataProcessor()
	
	sample := "  Hello   World!  This  is  a  test.  "
	
	cleaned := processor.CleanString(sample)
	println("Cleaned:", cleaned)
	
	upper := processor.NormalizeCase(sample, true)
	println("Uppercase:", upper)
	
	alnum := processor.ExtractAlphanumeric(sample)
	println("Alphanumeric only:", alnum)
}