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

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}

	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	for _, r := range data.Username {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' && r != '-' {
			return errors.New("username contains invalid characters")
		}
	}

	if !strings.Contains(data.Email, "@") || !strings.Contains(data.Email, ".") {
		return errors.New("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func TransformUserData(data UserData) (UserData, error) {
	if err := ValidateUserData(data); err != nil {
		return UserData{}, err
	}

	return UserData{
		Username: NormalizeUsername(data.Username),
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Age:      data.Age,
	}, nil
}