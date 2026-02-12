
package main

import (
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Username string
	Email    string
	Active   bool
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

func IsValidEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(email)
}

func ProcessUserInput(username, email string) (User, error) {
	if !ValidateUsername(username) {
		return User{}, ErrInvalidUsername
	}

	sanitizedEmail := SanitizeEmail(email)
	if !IsValidEmail(sanitizedEmail) {
		return User{}, ErrInvalidEmail
	}

	user := User{
		Username: username,
		Email:    sanitizedEmail,
		Active:   true,
	}

	return user, nil
}

var (
	ErrInvalidUsername = errors.New("invalid username format")
	ErrInvalidEmail    = errors.New("invalid email format")
)