
package main

import (
    "regexp"
    "strings"
)

type User struct {
    Username string
    Email    string
}

func ValidateUsername(username string) bool {
    if len(username) < 3 || len(username) > 20 {
        return false
    }
    validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    return validUsername.MatchString(username)
}

func SanitizeEmail(email string) string {
    trimmedEmail := strings.TrimSpace(email)
    return strings.ToLower(trimmedEmail)
}

func ValidateEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    return emailRegex.MatchString(email)
}

func ProcessUserInput(username, email string) (User, error) {
    if !ValidateUsername(username) {
        return User{}, ErrInvalidUsername
    }

    sanitizedEmail := SanitizeEmail(email)
    if !ValidateEmail(sanitizedEmail) {
        return User{}, ErrInvalidEmail
    }

    return User{
        Username: username,
        Email:    sanitizedEmail,
    }, nil
}

var (
    ErrInvalidUsername = errors.New("invalid username format")
    ErrInvalidEmail    = errors.New("invalid email format")
)