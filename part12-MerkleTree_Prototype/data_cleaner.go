
package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(input)
	
	// Remove extra internal whitespace
	re := regexp.MustCompile(`\s+`)
	normalized := re.ReplaceAllString(trimmed, " ")
	
	// Convert to lowercase
	lowercased := strings.ToLower(normalized)
	
	return lowercased
}

func RemoveSpecialChars(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, input)
}

func NormalizeWhitespace(input string) string {
	// Replace various whitespace characters with standard space
	re := regexp.MustCompile(`[\t\n\r\f\v]+`)
	return re.ReplaceAllString(input, " ")
}