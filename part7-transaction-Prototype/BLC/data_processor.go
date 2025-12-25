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

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}
	if data.Age < 18 || data.Age > 120 {
		return errors.New("age must be between 18 and 120")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
	transformedUsername := TransformUsername(username)
	userData := UserData{
		Username: transformedUsername,
		Email:    strings.TrimSpace(email),
		Age:      age,
	}
	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}
	return userData, nil
}
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

func ValidateJSONStructure(rawData string) (bool, error) {
    var data map[string]interface{}
    decoder := json.NewDecoder(strings.NewReader(rawData))
    decoder.DisallowUnknownFields()

    if err := decoder.Decode(&data); err != nil {
        return false, fmt.Errorf("invalid JSON structure: %w", err)
    }

    if len(data) == 0 {
        return false, fmt.Errorf("JSON data is empty")
    }

    for key, value := range data {
        if strings.TrimSpace(key) == "" {
            return false, fmt.Errorf("JSON contains empty key")
        }
        if value == nil {
            return false, fmt.Errorf("JSON key '%s' has nil value", key)
        }
    }

    return true, nil
}

func ExtractJSONKeys(rawData string) ([]string, error) {
    var data map[string]interface{}
    if err := json.Unmarshal([]byte(rawData), &data); err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
    }

    keys := make([]string, 0, len(data))
    for key := range data {
        keys = append(keys, key)
    }
    return keys, nil
}