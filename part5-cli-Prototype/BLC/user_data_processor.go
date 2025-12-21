package main

import (
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Username string
	Email    string
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validPattern.MatchString(username)
}

func SanitizeEmail(email string) string {
	trimmed := strings.TrimSpace(email)
	return strings.ToLower(trimmed)
}

func ValidateEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(email)
}

func ProcessUserInput(username, email string) (*User, error) {
	if !ValidateUsername(username) {
		return nil, &InputError{Field: "username", Message: "invalid username format"}
	}

	sanitizedEmail := SanitizeEmail(email)
	if !ValidateEmail(sanitizedEmail) {
		return nil, &InputError{Field: "email", Message: "invalid email address"}
	}

	return &User{
		Username: username,
		Email:    sanitizedEmail,
	}, nil
}

type InputError struct {
	Field   string
	Message string
}

func (e *InputError) Error() string {
	return e.Field + ": " + e.Message
}