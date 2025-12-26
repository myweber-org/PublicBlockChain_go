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

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUser(u User) (bool, []string) {
	var errors []string

	if u.Username == "" {
		errors = append(errors, "username cannot be empty")
	} else if len(u.Username) < 3 {
		errors = append(errors, "username must be at least 3 characters")
	}

	if u.Email == "" {
		errors = append(errors, "email cannot be empty")
	} else if !emailRegex.MatchString(u.Email) {
		errors = append(errors, "invalid email format")
	}

	return len(errors) == 0, errors
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func SanitizeEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	return email
}