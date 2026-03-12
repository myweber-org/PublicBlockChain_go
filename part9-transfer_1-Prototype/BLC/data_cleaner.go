package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(strings.TrimSpace(input), " ")

	// Remove non-printable characters
	cleaned = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, cleaned)

	return cleaned
}

func NormalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

func TruncateString(input string, maxLength int) string {
	if len(input) <= maxLength {
		return input
	}
	if maxLength < 3 {
		return input[:maxLength]
	}
	return input[:maxLength-3] + "..."
}