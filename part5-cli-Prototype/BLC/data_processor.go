package main

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func SanitizeInput(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func TransformToSlug(input string) string {
	slug := strings.ToLower(input)
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func main() {
	email := "test@example.com"
	fmt.Printf("Email %s valid: %v\n", email, ValidateEmail(email))

	input := "  HELLO World!  "
	fmt.Printf("Sanitized '%s': '%s'\n", input, SanitizeInput(input))

	title := "Go Programming 101: Best Practices"
	fmt.Printf("Slug for '%s': '%s'\n", title, TransformToSlug(title))
}