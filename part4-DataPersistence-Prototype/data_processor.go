package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	Username string
	Email    string
	Age      int
}

func ValidateProfile(profile UserProfile) error {
	if strings.TrimSpace(profile.Username) == "" {
		return errors.New("username cannot be empty")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if profile.Age < 0 || profile.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func NormalizeProfile(profile UserProfile) UserProfile {
	normalized := profile
	normalized.Username = strings.ToLower(strings.TrimSpace(profile.Username))
	normalized.Email = strings.ToLower(strings.TrimSpace(profile.Email))
	return normalized
}

func ProcessUserData(profile UserProfile) (UserProfile, error) {
	if err := ValidateProfile(profile); err != nil {
		return UserProfile{}, err
	}
	return NormalizeProfile(profile), nil
}