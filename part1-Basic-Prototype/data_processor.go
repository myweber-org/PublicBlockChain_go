
package main

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
}

func ValidateUsername(username string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]{3,20}$", username)
	return matched
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

func TransformProfile(profile UserProfile) (UserProfile, error) {
	if !ValidateUsername(profile.Username) {
		return profile, fmt.Errorf("invalid username format")
	}

	if !ValidateEmail(profile.Email) {
		return profile, fmt.Errorf("invalid email format")
	}

	if profile.Age < 0 || profile.Age > 150 {
		return profile, fmt.Errorf("age must be between 0 and 150")
	}

	profile.Username = strings.TrimSpace(profile.Username)
	profile.Email = strings.ToLower(strings.TrimSpace(profile.Email))

	return profile, nil
}

func ProcessUserData(jsonData []byte) (UserProfile, error) {
	var profile UserProfile
	err := json.Unmarshal(jsonData, &profile)
	if err != nil {
		return profile, fmt.Errorf("failed to parse JSON: %v", err)
	}

	transformedProfile, err := TransformProfile(profile)
	if err != nil {
		return profile, fmt.Errorf("validation failed: %v", err)
	}

	return transformedProfile, nil
}

func main() {
	jsonInput := `{"id":1,"username":"john_doe","email":"John@Example.COM","age":25,"active":true}`

	processedProfile, err := ProcessUserData([]byte(jsonInput))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	output, _ := json.MarshalIndent(processedProfile, "", "  ")
	fmt.Printf("Processed Profile:\n%s\n", output)
}