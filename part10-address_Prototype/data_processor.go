package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

func ValidateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ValidateAge(age int) error {
	if age < 0 || age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func ProcessUserData(data UserData) (UserData, error) {
	if err := ValidateEmail(data.Email); err != nil {
		return UserData{}, err
	}

	sanitizedUsername := SanitizeUsername(data.Username)
	if sanitizedUsername == "" {
		return UserData{}, errors.New("username cannot be empty")
	}

	if err := ValidateAge(data.Age); err != nil {
		return UserData{}, err
	}

	return UserData{
		Email:    data.Email,
		Username: sanitizedUsername,
		Age:      data.Age,
	}, nil
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
	Tags      []string `json:"tags"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func FilterInactiveUsers(users []UserProfile) []UserProfile {
	var activeUsers []UserProfile
	for _, user := range users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers
}

func TransformUserData(user UserProfile) map[string]interface{} {
	return map[string]interface{}{
		"user_id":   user.ID,
		"name":      NormalizeUsername(user.Username),
		"contact":   user.Email,
		"age_group": determineAgeGroup(user.Age),
		"status":    user.Active,
		"metadata":  user.Tags,
	}
}

func determineAgeGroup(age int) string {
	switch {
	case age < 18:
		return "minor"
	case age >= 18 && age <= 35:
		return "young_adult"
	case age > 35 && age <= 60:
		return "adult"
	default:
		return "senior"
	}
}

func ProcessUserData(jsonData string) ([]map[string]interface{}, error) {
	var users []UserProfile
	err := json.Unmarshal([]byte(jsonData), &users)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	var validUsers []UserProfile
	for _, user := range users {
		if ValidateEmail(user.Email) && user.Age > 0 {
			validUsers = append(validUsers, user)
		}
	}

	activeUsers := FilterInactiveUsers(validUsers)
	var transformedData []map[string]interface{}
	for _, user := range activeUsers {
		transformedData = append(transformedData, TransformUserData(user))
	}

	return transformedData, nil
}

func main() {
	sampleData := `[
		{"id":1,"username":"JohnDoe","email":"john@example.com","age":25,"active":true,"tags":["developer","gamer"]},
		{"id":2,"username":"JaneSmith","email":"invalid-email","age":30,"active":false,"tags":["designer"]},
		{"id":3,"username":"BobWilson","email":"bob@test.org","age":17,"active":true,"tags":["student"]}
	]`

	result, err := ProcessUserData(sampleData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("Processed user data:")
	fmt.Println(string(output))
}