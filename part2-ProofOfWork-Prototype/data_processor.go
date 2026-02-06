package main

import (
	"errors"
	"strings"
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
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(data UserData) UserData {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
	return data
}

func ProcessUserInput(rawUsername string, rawEmail string, rawAge int) (UserData, error) {
	userData := UserData{
		Username: rawUsername,
		Email:    rawEmail,
		Age:      rawAge,
	}

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	userData = TransformUsername(userData)
	return userData, nil
}
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dp *DataProcessor) NormalizeEmail(email string) (string, bool) {
	cleaned := strings.ToLower(dp.CleanString(email))
	if dp.emailRegex.MatchString(cleaned) {
		return cleaned, true
	}
	return cleaned, false
}

func (dp *DataProcessor) ValidateAndProcess(input string) (string, []string) {
	cleaned := dp.CleanString(input)
	parts := strings.Fields(cleaned)
	
	var validEmails []string
	var otherParts []string
	
	for _, part := range parts {
		if normalized, valid := dp.NormalizeEmail(part); valid {
			validEmails = append(validEmails, normalized)
		} else {
			otherParts = append(otherParts, part)
		}
	}
	
	return strings.Join(otherParts, " "), validEmails
}