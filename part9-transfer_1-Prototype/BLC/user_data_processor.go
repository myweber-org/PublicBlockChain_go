package main

import (
	"regexp"
	"strings"
)

type User struct {
	Username string
	Email    string
	Age      int
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validUsername.MatchString(username)
}

func SanitizeEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	return email
}

func ValidateUserAge(age int) bool {
	return age >= 18 && age <= 120
}

func ProcessUserInput(username, email string, age int) (User, error) {
	if !ValidateUsername(username) {
		return User{}, ErrInvalidUsername
	}

	sanitizedEmail := SanitizeEmail(email)
	if !strings.Contains(sanitizedEmail, "@") {
		return User{}, ErrInvalidEmail
	}

	if !ValidateUserAge(age) {
		return User{}, ErrInvalidAge
	}

	return User{
		Username: username,
		Email:    sanitizedEmail,
		Age:      age,
	}, nil
}

var (
	ErrInvalidUsername = errors.New("invalid username format")
	ErrInvalidEmail    = errors.New("invalid email address")
	ErrInvalidAge      = errors.New("age must be between 18 and 120")
)