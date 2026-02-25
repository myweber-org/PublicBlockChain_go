
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	Email    string
	Username string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateProfile(p UserProfile) error {
	if !emailRegex.MatchString(p.Email) {
		return errors.New("invalid email format")
	}
	if strings.TrimSpace(p.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if p.Age < 0 || p.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(p *UserProfile) {
	p.Username = strings.ToLower(strings.TrimSpace(p.Username))
}

func ProcessUserProfile(p UserProfile) (UserProfile, error) {
	if err := ValidateProfile(p); err != nil {
		return p, err
	}
	TransformUsername(&p)
	return p, nil
}