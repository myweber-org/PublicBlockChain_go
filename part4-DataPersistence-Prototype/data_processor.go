
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

func (dp *DataProcessor) ExtractDomain(email string) (string, bool) {
	if !dp.ValidateEmail(email) {
		return "", false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", false
	}
	return parts[1], true
}

func (dp *DataProcessor) NormalizeWhitespace(input string) string {
	return dp.whitespaceRegex.ReplaceAllString(input, " ")
}package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	emailRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return &DataProcessor{
		emailRegex: regexp.MustCompile(pattern),
	}
}

func (dp *DataProcessor) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToLower(trimmed)
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, bool) {
	sanitizedName := dp.SanitizeInput(name)
	sanitizedEmail := dp.SanitizeInput(email)

	if sanitizedName == "" || sanitizedEmail == "" {
		return "", false
	}

	if !dp.ValidateEmail(sanitizedEmail) {
		return "", false
	}

	return sanitizedName + " <" + sanitizedEmail + ">", true
}