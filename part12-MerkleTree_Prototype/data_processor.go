package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	if !validPattern.MatchString(trimmed) {
		return "", false
	}
	return trimmed, true
}

func ValidateEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(email)
}

func ProcessUserData(username, email string) (map[string]interface{}, error) {
	sanitizedUsername, ok := SanitizeUsername(username)
	if !ok {
		return nil, &InvalidDataError{Field: "username", Value: username}
	}

	if !ValidateEmail(email) {
		return nil, &InvalidDataError{Field: "email", Value: email}
	}

	return map[string]interface{}{
		"username": sanitizedUsername,
		"email":    strings.ToLower(email),
		"status":   "processed",
	}, nil
}

type InvalidDataError struct {
	Field string
	Value string
}

func (e *InvalidDataError) Error() string {
	return "invalid data for field: " + e.Field
}