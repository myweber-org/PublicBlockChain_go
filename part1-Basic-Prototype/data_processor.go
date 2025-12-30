
package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int
	Email     string
	Username  string
	BirthDate string
	CreatedAt time.Time
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return errors.New("invalid user ID")
	}

	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	if strings.ContainsAny(profile.Username, "!@#$%^&*()") {
		return errors.New("username contains invalid characters")
	}

	_, err := time.Parse("2006-01-02", profile.BirthDate)
	if err != nil {
		return errors.New("invalid birth date format, expected YYYY-MM-DD")
	}

	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func CalculateAge(birthDate string) (int, error) {
	birth, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	age := now.Year() - birth.Year()

	if now.YearDay() < birth.YearDay() {
		age--
	}

	return age, nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateUserProfile(profile); err != nil {
		return profile, err
	}

	profile.Username = TransformUsername(profile.Username)
	return profile, nil
}