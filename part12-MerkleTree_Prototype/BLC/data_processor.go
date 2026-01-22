
package data_processor

import (
	"encoding/json"
	"fmt"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

type UserData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

func ParseAndValidateUserData(rawData []byte) (*UserData, error) {
	var user UserData
	if err := json.Unmarshal(rawData, &user); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if user.ID <= 0 {
		return nil, ValidationError{Field: "id", Message: "must be positive integer"}
	}

	if user.Username == "" {
		return nil, ValidationError{Field: "username", Message: "cannot be empty"}
	}

	if !isValidEmail(user.Email) {
		return nil, ValidationError{Field: "email", Message: "invalid email format"}
	}

	return &user, nil
}

func isValidEmail(email string) bool {
	const minEmailLength = 5
	if len(email) < minEmailLength {
		return false
	}

	hasAt := false
	hasDot := false
	for i, char := range email {
		if char == '@' {
			if hasAt || i == 0 || i == len(email)-1 {
				return false
			}
			hasAt = true
		}
		if char == '.' && hasAt && i > 0 && i < len(email)-1 {
			hasDot = true
		}
	}
	return hasAt && hasDot
}