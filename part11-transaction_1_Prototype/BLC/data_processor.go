package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Comments string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	re := regexp.MustCompile(`[<>"'&]`)
	return re.ReplaceAllString(input, "")
}

func ValidateUserData(data UserData) (bool, []string) {
	var errors []string

	if len(data.Username) < 3 || len(data.Username) > 20 {
		errors = append(errors, "Username must be between 3 and 20 characters")
	}

	if !emailRegex.MatchString(data.Email) {
		errors = append(errors, "Invalid email format")
	}

	if len(data.Comments) > 500 {
		errors = append(errors, "Comments cannot exceed 500 characters")
	}

	return len(errors) == 0, errors
}

func ProcessUserData(data UserData) UserData {
	return UserData{
		Username: SanitizeInput(data.Username),
		Email:    strings.ToLower(SanitizeInput(data.Email)),
		Comments: SanitizeInput(data.Comments),
	}
}package main

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
	if strings.TrimSpace(data.Username) == "" {
		return UserData{}, fmt.Errorf("username cannot be empty")
	}

	if !strings.Contains(data.Email, "@") {
		return UserData{}, fmt.Errorf("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return UserData{}, fmt.Errorf("age must be between 0 and 150")
	}

	transformed := UserData{
		Username: strings.ToLower(strings.TrimSpace(data.Username)),
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Age:      data.Age,
	}

	return transformed, nil
}

func main() {
	sampleData := UserData{
		Username: "  TestUser  ",
		Email:    "TEST@EXAMPLE.COM",
		Age:      25,
	}

	result, err := ValidateAndTransform(sampleData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Original: %+v\n", sampleData)
	fmt.Printf("Processed: %+v\n", result)
}