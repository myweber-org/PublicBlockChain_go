package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Password string
}

func ValidateUserData(data UserData) (bool, []string) {
	var errors []string

	if len(data.Username) < 3 || len(data.Username) > 20 {
		errors = append(errors, "Username must be between 3 and 20 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(data.Email) {
		errors = append(errors, "Invalid email format")
	}

	if len(data.Password) < 8 {
		errors = append(errors, "Password must be at least 8 characters")
	}

	return len(errors) == 0, errors
}

func SanitizeUserInput(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	return input
}

func ProcessUserRegistration(data UserData) (bool, []string) {
	data.Username = SanitizeUserInput(data.Username)
	data.Email = SanitizeUserInput(data.Email)

	return ValidateUserData(data)
}