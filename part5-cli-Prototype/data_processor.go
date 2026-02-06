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

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return fmt.Errorf("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(data *UserData) {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
	user := UserData{
		Username: username,
		Email:    email,
		Age:      age,
	}

	TransformUsername(&user)

	if err := ValidateUserData(user); err != nil {
		return UserData{}, err
	}

	return user, nil
}

func main() {
	user, err := ProcessUserInput("  JohnDoe  ", "john@example.com", 30)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	ID    int
	Name  string
	Email string
	Age   int
}

func ValidateUserData(data UserData) error {
	if data.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", data.ID)
	}
	if strings.TrimSpace(data.Name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return fmt.Errorf("invalid email format: %s", data.Email)
	}
	if data.Age < 0 || data.Age > 150 {
		return fmt.Errorf("age out of valid range: %d", data.Age)
	}
	return nil
}

func NormalizeUserData(data UserData) UserData {
	return UserData{
		ID:    data.ID,
		Name:  strings.TrimSpace(data.Name),
		Email: strings.ToLower(strings.TrimSpace(data.Email)),
		Age:   data.Age,
	}
}

func ProcessUserInput(rawData UserData) (UserData, error) {
	normalizedData := NormalizeUserData(rawData)
	if err := ValidateUserData(normalizedData); err != nil {
		return UserData{}, err
	}
	return normalizedData, nil
}

func main() {
	testData := UserData{
		ID:    1001,
		Name:  "  John Doe  ",
		Email: "  John@Example.COM  ",
		Age:   30,
	}

	processedData, err := ProcessUserInput(testData)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	fmt.Printf("Processed user data: %+v\n", processedData)
}