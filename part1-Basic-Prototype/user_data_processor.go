
package main

import (
    "regexp"
    "strings"
)

type User struct {
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

func ValidateUserAge(age int) bool {
    return age >= 13 && age <= 120
}

func ProcessUserInput(username, email string, age int) (User, error) {
    if !ValidateUsername(username) {
        return User{}, ErrInvalidUsername
    }

    sanitizedEmail := SanitizeEmail(email)
    if !ValidateUserAge(age) {
        return User{}, ErrInvalidAge
    }

    return User{
        Username: username,
        Email:    sanitizedEmail,
        Age:      age,
    }, nil
}

var (
    ErrInvalidUsername = errors.New("invalid username format")
    ErrInvalidAge      = errors.New("age must be between 13 and 120")
)