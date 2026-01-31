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

func ValidateAndTransform(data UserData) (UserData, error) {
	var processed UserData

	if data.Username == "" {
		return processed, fmt.Errorf("username cannot be empty")
	}
	processed.Username = strings.TrimSpace(data.Username)

	if !strings.Contains(data.Email, "@") {
		return processed, fmt.Errorf("invalid email format")
	}
	processed.Email = strings.ToLower(strings.TrimSpace(data.Email))

	if data.Age < 0 || data.Age > 150 {
		return processed, fmt.Errorf("age must be between 0 and 150")
	}
	processed.Age = data.Age

	return processed, nil
}

func main() {
	sample := UserData{
		Username: "  JohnDoe  ",
		Email:    "John@Example.COM",
		Age:      25,
	}

	result, err := ValidateAndTransform(sample)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Processed: %+v\n", result)
}