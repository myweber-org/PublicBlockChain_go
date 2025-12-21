package main

import (
	"regexp"
	"strings"
)

type UserData struct {
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
	trimmed := strings.TrimSpace(email)
	return strings.ToLower(trimmed)
}

func ValidateAge(age int) bool {
	return age >= 18 && age <= 120
}

func ProcessUserData(data UserData) (UserData, error) {
	if !ValidateUsername(data.Username) {
		return UserData{}, ErrInvalidUsername
	}

	data.Email = SanitizeEmail(data.Email)

	if !ValidateAge(data.Age) {
		return UserData{}, ErrInvalidAge
	}

	return data, nil
}

var (
	ErrInvalidUsername = errors.New("invalid username format")
	ErrInvalidAge      = errors.New("age must be between 18 and 120")
)