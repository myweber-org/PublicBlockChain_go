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
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID    string
	Name  string
	Email string
	Valid bool
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		if len(line) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(line[0]),
			Name:  strings.TrimSpace(line[1]),
			Email: strings.TrimSpace(line[2]),
			Valid: validateRecord(strings.TrimSpace(line[0]), strings.TrimSpace(line[2])),
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(id, email string) bool {
	if id == "" || email == "" {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func GenerateReport(records []DataRecord) {
	fmt.Printf("Total records processed: %d\n", len(records))
	validCount := 0
	for _, record := range records {
		if record.Valid {
			validCount++
		}
	}
	fmt.Printf("Valid records: %d\n", validCount)
	fmt.Printf("Invalid records: %d\n", len(records)-validCount)
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
	Timestamp string `json:"timestamp"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func TransformProfile(profile UserProfile) (UserProfile, error) {
	if profile.Age < 0 || profile.Age > 150 {
		return profile, fmt.Errorf("invalid age: %d", profile.Age)
	}

	if !ValidateEmail(profile.Email) {
		return profile, fmt.Errorf("invalid email format: %s", profile.Email)
	}

	profile.Username = SanitizeUsername(profile.Username)

	if profile.Timestamp == "" {
		profile.Timestamp = "2024-01-01T00:00:00Z"
	}

	return profile, nil
}

func ProcessUserData(jsonData []byte) ([]UserProfile, error) {
	var profiles []UserProfile
	err := json.Unmarshal(jsonData, &profiles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	var validProfiles []UserProfile
	for _, profile := range profiles {
		transformed, err := TransformProfile(profile)
		if err != nil {
			fmt.Printf("Skipping profile ID %d: %v\n", profile.ID, err)
			continue
		}
		validProfiles = append(validProfiles, transformed)
	}

	return validProfiles, nil
}

func main() {
	jsonInput := `[
		{"id":1,"username":"  JohnDoe  ","email":"john@example.com","age":25,"active":true,"timestamp":"2024-01-15T10:30:00Z"},
		{"id":2,"username":"JaneSmith","email":"invalid-email","age":-5,"active":false,"timestamp":""},
		{"id":3,"username":"BOB_LEE","email":"bob@company.org","age":30,"active":true,"timestamp":"2024-01-16T14:45:00Z"}
	]`

	processed, err := ProcessUserData([]byte(jsonInput))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Printf("Successfully processed %d profiles\n", len(processed))
	for _, p := range processed {
		fmt.Printf("ID: %d, Username: %s, Email: %s, Age: %d\n", p.ID, p.Username, p.Email, p.Age)
	}
}