
package data

import (
	"errors"
	"regexp"
	"strings"
)

type Record struct {
	ID    string
	Email string
	Score int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateRecord(r Record) error {
	if r.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}
	if r.Score < 0 || r.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func TransformRecords(records []Record) ([]Record, error) {
	var processed []Record
	for _, r := range records {
		if err := ValidateRecord(r); err != nil {
			return nil, err
		}
		r.Email = NormalizeEmail(r.Email)
		processed = append(processed, r)
	}
	return processed, nil
}

func CalculateAverage(records []Record) float64 {
	if len(records) == 0 {
		return 0.0
	}
	total := 0
	for _, r := range records {
		total += r.Score
	}
	return float64(total) / float64(len(records))
}
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

func processCSVFile(inputPath string, outputPath string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    reader := csv.NewReader(inputFile)
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    lineNumber := 0
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV at line %d: %w", lineNumber, err)
        }

        lineNumber++
        if lineNumber == 1 {
            if err := writer.Write(record); err != nil {
                return fmt.Errorf("error writing header: %w", err)
            }
            continue
        }

        cleanedRecord := cleanRecord(record)
        if cleanedRecord == nil {
            continue
        }

        if err := writer.Write(cleanedRecord); err != nil {
            return fmt.Errorf("error writing record at line %d: %w", lineNumber, err)
        }
    }

    return nil
}

func cleanRecord(record []string) []string {
    cleaned := make([]string, len(record))
    for i, field := range record {
        cleaned[i] = strings.TrimSpace(field)
        if cleaned[i] == "" {
            return nil
        }
    }
    return cleaned
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_processor <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := processCSVFile(inputFile, outputFile); err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %s to %s\n", inputFile, outputFile)
}
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	Username string
	Email    string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateProfile(profile UserProfile) error {
	if strings.TrimSpace(profile.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}
	if profile.Age < 0 || profile.Age > 120 {
		return errors.New("age must be between 0 and 120")
	}
	return nil
}

func TransformUsername(profile *UserProfile) {
	profile.Username = strings.ToLower(strings.TrimSpace(profile.Username))
}

func ProcessUserData(profile UserProfile) (UserProfile, error) {
	TransformUsername(&profile)
	if err := ValidateProfile(profile); err != nil {
		return UserProfile{}, err
	}
	return profile, nil
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

func TransformUserData(users []UserProfile) ([]map[string]interface{}, error) {
	var transformed []map[string]interface{}
	
	for _, user := range users {
		if !ValidateEmail(user.Email) {
			return nil, fmt.Errorf("invalid email for user %d", user.ID)
		}
		
		data := map[string]interface{}{
			"user_id":   user.ID,
			"username":  NormalizeUsername(user.Username),
			"email":     user.Email,
			"age_group": categorizeAge(user.Age),
			"tag_count": len(user.Tags),
			"metadata": map[string]interface{}{
				"active": user.Active,
				"tags":   user.Tags,
			},
		}
		transformed = append(transformed, data)
	}
	
	return transformed, nil
}

func categorizeAge(age int) string {
	switch {
	case age < 18:
		return "minor"
	case age >= 18 && age < 65:
		return "adult"
	default:
		return "senior"
	}
}

func ProcessUserJSON(jsonData []byte) (string, error) {
	var users []UserProfile
	err := json.Unmarshal(jsonData, &users)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}
	
	activeUsers := FilterInactiveUsers(users)
	transformed, err := TransformUserData(activeUsers)
	if err != nil {
		return "", fmt.Errorf("transformation failed: %v", err)
	}
	
	result, err := json.MarshalIndent(transformed, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}
	
	return string(result), nil
}

func main() {
	sampleJSON := `[
		{"id":1,"username":"JohnDoe","email":"john@example.com","age":25,"active":true,"tags":["golang","backend"]},
		{"id":2,"username":"JaneSmith","email":"jane@example.org","age":17,"active":false,"tags":["frontend"]},
		{"id":3,"username":"BobWilson","email":"bob@test.com","age":70,"active":true,"tags":["devops","cloud"]}
	]`
	
	result, err := ProcessUserJSON([]byte(sampleJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("Processed user data:")
	fmt.Println(result)
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
	ID      string
	Name    string
	Value   string
	IsValid bool
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Name:  strings.TrimSpace(row[1]),
			Value: strings.TrimSpace(row[2]),
		}

		record.IsValid = validateRecord(record)
		records = append(records, record)
	}

	return records, nil
}

func validateRecord(record DataRecord) bool {
	if record.ID == "" || record.Name == "" {
		return false
	}

	if len(record.Value) > 100 {
		return false
	}

	return true
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if record.IsValid {
			valid = append(valid, record)
		}
	}
	return valid
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	validRecords := FilterValidRecords(records)
	fmt.Printf("Processed %d records, %d valid\n", len(records), len(validRecords))

	for _, record := range validRecords {
		fmt.Printf("ID: %s, Name: %s\n", record.ID, record.Name)
	}
}