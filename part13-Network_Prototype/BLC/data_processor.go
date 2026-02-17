package main

import (
	"errors"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

func ValidateUserData(data UserData) error {
	if data.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	if data.Username == "" {
		return errors.New("username is required")
	}
	if len(data.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	transformedUsername := TransformUsername(username)
	userData := UserData{
		Email:    strings.TrimSpace(email),
		Username: transformedUsername,
		Age:      age,
	}
	err := ValidateUserData(userData)
	return userData, err
}