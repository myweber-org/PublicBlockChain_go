
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
}