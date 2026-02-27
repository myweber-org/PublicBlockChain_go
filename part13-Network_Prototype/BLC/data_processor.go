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
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return errors.New("username can only contain letters, numbers and underscores")
	}
	return nil
}

func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ProcessUserData(data UserData) (UserData, error) {
	if err := ValidateUsername(data.Username); err != nil {
		return UserData{}, err
	}
	
	if err := ValidateEmail(data.Email); err != nil {
		return UserData{}, err
	}
	
	data.Email = NormalizeEmail(data.Email)
	
	if data.Age < 0 || data.Age > 120 {
		return UserData{}, errors.New("age must be between 0 and 120")
	}
	
	return data, nil
}

func TransformUsername(username string) string {
	return strings.ReplaceAll(username, "_", "-")
}