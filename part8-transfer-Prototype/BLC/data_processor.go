
package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", ErrEmptyInput
	}

	pattern := `^[a-zA-Z0-9_\-\.]+$`
	matched, err := regexp.MatchString(pattern, trimmed)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", ErrInvalidCharacters
	}

	if len(trimmed) > 50 {
		return "", ErrInputTooLong
	}
	return trimmed, nil
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

var (
	ErrEmptyInput        = errors.New("input cannot be empty")
	ErrInvalidCharacters = errors.New("input contains invalid characters")
	ErrInputTooLong      = errors.New("input exceeds maximum length")
)