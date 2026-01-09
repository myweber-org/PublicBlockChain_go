package main

import (
	"regexp"
	"strings"
)

// SanitizeUserInput removes potentially dangerous HTML/script content
func SanitizeUserInput(input string) string {
	// Remove script tags and their content
	scriptRegex := regexp.MustCompile(`<script\b[^>]*>(.*?)</script>`)
	sanitized := scriptRegex.ReplaceAllString(input, "")

	// Remove common event handlers
	eventHandlers := []string{
		`onclick`, `onload`, `onerror`, `onmouseover`,
		`onkeypress`, `onsubmit`, `onfocus`, `onblur`,
	}
	for _, handler := range eventHandlers {
		pattern := regexp.MustCompile(`\s*` + regexp.QuoteMeta(handler) + `\s*=\s*["'][^"']*["']`)
		sanitized = pattern.ReplaceAllString(sanitized, "")
	}

	// Escape remaining HTML special characters
	sanitized = strings.ReplaceAll(sanitized, "&", "&amp;")
	sanitized = strings.ReplaceAll(sanitized, "<", "&lt;")
	sanitized = strings.ReplaceAll(sanitized, ">", "&gt;")
	sanitized = strings.ReplaceAll(sanitized, "\"", "&quot;")
	sanitized = strings.ReplaceAll(sanitized, "'", "&#39;")

	return strings.TrimSpace(sanitized)
}

// ValidateEmail checks if a string is a valid email address
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func main() {
	// Example usage
	testInput := `<script>alert("xss")</script><img src=x onerror=alert(1)>`
	cleaned := SanitizeUserInput(testInput)
	println("Sanitized output:", cleaned)

	testEmail := "user@example.com"
	println("Email valid:", ValidateEmail(testEmail))
}