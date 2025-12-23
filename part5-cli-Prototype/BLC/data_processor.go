
package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processUserData(data UserData) (UserData, error) {
	if data.Username == "" {
		return data, fmt.Errorf("username cannot be empty")
	}
	data.Username = normalizeUsername(data.Username)

	if !validateEmail(data.Email) {
		return data, fmt.Errorf("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return data, fmt.Errorf("age must be between 0 and 150")
	}

	return data, nil
}

func main() {
	user := UserData{
		Username: "  JohnDoe  ",
		Email:    "john@example.com",
		Age:      30,
	}

	processedUser, err := processUserData(user)
	if err != nil {
		fmt.Printf("Error processing user data: %v\n", err)
		return
	}

	fmt.Printf("Processed user: %+v\n", processedUser)
}package main

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
	if len(data.Username) < 3 || len(data.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
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
	err := ValidateUserData(userData)
	if err != nil {
		return UserData{}, err
	}
	return userData, nil
}