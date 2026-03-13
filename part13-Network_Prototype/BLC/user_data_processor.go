package main

import (
	"errors"
	"strings"
	"unicode"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(username) > 20 {
		return errors.New("username must not exceed 20 characters")
	}
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return errors.New("username can only contain letters, digits, underscores and hyphens")
		}
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func IsValidAge(age int) bool {
	return age >= 0 && age <= 150
}

func ProcessUserInput(username, email string, age int) (*UserData, error) {
	if err := ValidateUsername(username); err != nil {
		return nil, err
	}

	normalizedEmail := NormalizeEmail(email)
	if !strings.Contains(normalizedEmail, "@") {
		return nil, errors.New("invalid email format")
	}

	if !IsValidAge(age) {
		return nil, errors.New("age must be between 0 and 150")
	}

	return &UserData{
		Username: username,
		Email:    normalizedEmail,
		Age:      age,
	}, nil
}