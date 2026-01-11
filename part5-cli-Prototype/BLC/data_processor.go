package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

func ValidateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ProcessUserData(data UserData) (UserData, error) {
	if err := ValidateEmail(data.Email); err != nil {
		return UserData{}, err
	}

	data.Username = SanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return UserData{}, errors.New("age must be between 0 and 150")
	}

	return data, nil
}