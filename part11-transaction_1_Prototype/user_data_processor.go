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
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]{3,20}$", username)
	return matched
}

func SanitizeEmail(email string) string {
	trimmed := strings.TrimSpace(email)
	return strings.ToLower(trimmed)
}

func ValidateAge(age int) bool {
	return age >= 13 && age <= 120
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
	ErrInvalidUsername = NewValidationError("invalid username format")
	ErrInvalidAge      = NewValidationError("age must be between 13 and 120")
)

type ValidationError struct {
	Message string
}

func NewValidationError(msg string) ValidationError {
	return ValidationError{Message: msg}
}

func (e ValidationError) Error() string {
	return e.Message
}