package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Comments string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	re := regexp.MustCompile(`[<>"'&]`)
	return re.ReplaceAllString(input, "")
}

func ValidateUserData(data UserData) (bool, []string) {
	var errors []string

	if len(data.Username) < 3 || len(data.Username) > 20 {
		errors = append(errors, "Username must be between 3 and 20 characters")
	}

	if !emailRegex.MatchString(data.Email) {
		errors = append(errors, "Invalid email format")
	}

	if len(data.Comments) > 500 {
		errors = append(errors, "Comments cannot exceed 500 characters")
	}

	return len(errors) == 0, errors
}

func ProcessUserData(data UserData) UserData {
	return UserData{
		Username: SanitizeInput(data.Username),
		Email:    strings.ToLower(SanitizeInput(data.Email)),
		Comments: SanitizeInput(data.Comments),
	}
}