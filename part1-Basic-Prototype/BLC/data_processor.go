
package main

import (
	"errors"
	"strings"
	"unicode"
)

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return errors.New("username can only contain letters, digits, underscores, and hyphens")
		}
	}

	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func SanitizeInput(input string) string {
	replacer := strings.NewReplacer(
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(input)
}

func ValidatePasswordStrength(password string) (bool, []string) {
	var issues []string

	if len(password) < 8 {
		issues = append(issues, "password must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if !hasUpper {
		issues = append(issues, "password must contain at least one uppercase letter")
	}
	if !hasLower {
		issues = append(issues, "password must contain at least one lowercase letter")
	}
	if !hasDigit {
		issues = append(issues, "password must contain at least one digit")
	}
	if !hasSpecial {
		issues = append(issues, "password must contain at least one special character")
	}

	return len(issues) == 0, issues
}
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        string
	Email     string
	Username  string
	Age       int
	Active    bool
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateProfile(p UserProfile) error {
	if p.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if !emailRegex.MatchString(p.Email) {
		return errors.New("invalid email format")
	}
	if len(p.Username) < 3 || len(p.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if p.Age < 0 || p.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformProfile(p UserProfile) UserProfile {
	p.Username = strings.ToLower(strings.TrimSpace(p.Username))
	p.Email = strings.ToLower(strings.TrimSpace(p.Email))
	return p
}

func ProcessUserProfile(p UserProfile) (UserProfile, error) {
	if err := ValidateProfile(p); err != nil {
		return UserProfile{}, err
	}
	transformed := TransformProfile(p)
	return transformed, nil
}