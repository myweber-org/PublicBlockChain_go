
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

func ValidateEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    return emailRegex.MatchString(email)
}

func ProcessUserInput(username, email string) (*User, error) {
    if !ValidateUsername(username) {
        return nil, fmt.Errorf("invalid username format")
    }
    
    sanitizedEmail := SanitizeEmail(email)
    if !ValidateEmail(sanitizedEmail) {
        return nil, fmt.Errorf("invalid email format")
    }
    
    return &User{
        Username: username,
        Email:    sanitizedEmail,
    }, nil
}