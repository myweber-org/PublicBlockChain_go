
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

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ExtractDomain(email string) (string, bool) {
	if !dp.ValidateEmail(email) {
		return "", false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", false
	}
	return parts[1], true
}

func (dp *DataProcessor) NormalizeWhitespace(input string) string {
	return dp.whitespaceRegex.ReplaceAllString(input, " ")
}
package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput removes potentially harmful characters and trims whitespace.
// It returns an empty string if the input contains disallowed patterns.
func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return ""
	}

	// Reject input containing script tags or SQL comment patterns
	maliciousPattern := regexp.MustCompile(`(?i)<script|--|/\*|\*/`)
	if maliciousPattern.MatchString(trimmed) {
		return ""
	}

	// Allow only alphanumeric, spaces, and basic punctuation
	safePattern := regexp.MustCompile(`[^a-zA-Z0-9\s.,!?-]`)
	sanitized := safePattern.ReplaceAllString(trimmed, "")

	return strings.TrimSpace(sanitized)
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

func ValidateUserProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return fmt.Errorf("invalid ID: must be positive integer")
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	if !usernameRegex.MatchString(profile.Username) {
		return fmt.Errorf("invalid username: must be 3-20 alphanumeric characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return fmt.Errorf("invalid email format")
	}

	if profile.Age < 0 || profile.Age > 150 {
		return fmt.Errorf("invalid age: must be between 0 and 150")
	}

	return nil
}

func TransformProfile(profile UserProfile) UserProfile {
	transformed := profile
	transformed.Username = strings.ToLower(transformed.Username)
	transformed.Email = strings.ToLower(transformed.Email)
	
	uniqueTags := make(map[string]bool)
	var cleanedTags []string
	for _, tag := range transformed.Tags {
		cleanTag := strings.TrimSpace(tag)
		if cleanTag != "" && !uniqueTags[cleanTag] {
			uniqueTags[cleanTag] = true
			cleanedTags = append(cleanedTags, cleanTag)
		}
	}
	transformed.Tags = cleanedTags
	
	return transformed
}

func ProcessUserProfile(data []byte) (UserProfile, error) {
	var profile UserProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return UserProfile{}, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if err := ValidateUserProfile(profile); err != nil {
		return UserProfile{}, fmt.Errorf("validation failed: %v", err)
	}

	transformedProfile := TransformProfile(profile)
	return transformedProfile, nil
}

func main() {
	jsonData := []byte(`{
		"id": 123,
		"username": "TestUser_123",
		"email": "TEST@EXAMPLE.COM",
		"age": 30,
		"active": true,
		"tags": ["go", " backend", "go ", "dev", ""]
	}`)

	processedProfile, err := ProcessUserProfile(jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Processed Profile: %+v\n", processedProfile)
}