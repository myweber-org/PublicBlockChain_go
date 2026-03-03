
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
}package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	emailRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return &DataProcessor{
		emailRegex: regexp.MustCompile(pattern),
	}
}

func (dp *DataProcessor) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToLower(trimmed)
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, bool) {
	sanitizedName := dp.SanitizeInput(name)
	sanitizedEmail := dp.SanitizeInput(email)

	if sanitizedName == "" || sanitizedEmail == "" {
		return "", false
	}

	if !dp.ValidateEmail(sanitizedEmail) {
		return "", false
	}

	return sanitizedName + " <" + sanitizedEmail + ">", true
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
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func FilterInactiveUsers(users []UserProfile) []UserProfile {
	var activeUsers []UserProfile
	for _, user := range users {
		if user.Active && user.Age >= 18 {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers
}

func ProcessUserData(inputJSON string) ([]UserProfile, error) {
	var users []UserProfile
	err := json.Unmarshal([]byte(inputJSON), &users)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	for i := range users {
		users[i].Username = NormalizeUsername(users[i].Username)
		
		if !ValidateEmail(users[i].Email) {
			return nil, fmt.Errorf("invalid email for user %d", users[i].ID)
		}
	}

	return FilterInactiveUsers(users), nil
}

func main() {
	jsonData := `[
		{"id":1,"username":" JohnDoe ","email":"john@example.com","age":25,"active":true,"tags":["golang","backend"]},
		{"id":2,"username":"jane_smith","email":"invalid-email","age":30,"active":true,"tags":["frontend"]},
		{"id":3,"username":"inactive_user","email":"test@domain.com","age":16,"active":false,"tags":[]}
	]`

	processedUsers, err := ProcessUserData(jsonData)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Printf("Valid active users: %d\n", len(processedUsers))
	for _, user := range processedUsers {
		fmt.Printf("ID: %d, Username: %s, Email: %s\n", user.ID, user.Username, user.Email)
	}
}