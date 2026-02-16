
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
    Email   string
    Active  string
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
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

        if lineNumber == 1 {
            continue
        }

        if len(row) < 4 {
            return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
        }

        record := DataRecord{
            ID:     strings.TrimSpace(row[0]),
            Name:   strings.TrimSpace(row[1]),
            Email:  strings.TrimSpace(row[2]),
            Active: strings.TrimSpace(row[3]),
        }

        if record.ID == "" || record.Name == "" {
            return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
        }

        records = append(records, record)
    }

    return records, nil
}

func ValidateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FilterActiveUsers(records []DataRecord) []DataRecord {
    var activeUsers []DataRecord
    for _, record := range records {
        if record.Active == "true" && ValidateEmail(record.Email) {
            activeUsers = append(activeUsers, record)
        }
    }
    return activeUsers
}

func GenerateReport(records []DataRecord) {
    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Println("Active users with valid emails:")
    for _, user := range FilterActiveUsers(records) {
        fmt.Printf("  - %s (%s)\n", user.Name, user.Email)
    }
}
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func sanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func transformUserData(rawData []byte) (*UserData, error) {
	var user UserData
	err := json.Unmarshal(rawData, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if !validateEmail(user.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	user.Username = sanitizeUsername(user.Username)

	if user.Age < 0 || user.Age > 150 {
		return nil, fmt.Errorf("age out of valid range")
	}

	return &user, nil
}

func main() {
	jsonData := []byte(`{"email":"test@example.com","username":"  john_doe  ","age":25}`)
	processedUser, err := transformUserData(jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", processedUser)
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func ProcessUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if !ValidateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range")
	}

	data.Username = SanitizeUsername(data.Username)

	return &data, nil
}

func main() {
	jsonData := []byte(`{"email":"test@example.com","username":"  JohnDoe  ","age":25}`)
	processed, err := ProcessUserData(jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed: %+v\n", processed)
}
package data

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
	ErrEmptyData    = errors.New("data cannot be empty")
)

type DataRecord struct {
	ID        string
	Timestamp time.Time
	Value     float64
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return ErrInvalidInput
	}
	if record.Value < 0 {
		return errors.New("value cannot be negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp must be set")
	}
	return nil
}

func NormalizeString(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func FilterRecords(records []DataRecord, minValue float64) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Value >= minValue {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func CalculateAverage(records []DataRecord) (float64, error) {
	if len(records) == 0 {
		return 0, ErrEmptyData
	}
	
	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	return sum / float64(len(records)), nil
}

func MergeTags(records []DataRecord) []string {
	tagMap := make(map[string]bool)
	for _, record := range records {
		for _, tag := range record.Tags {
			normalized := NormalizeString(tag)
			if normalized != "" {
				tagMap[normalized] = true
			}
		}
	}
	
	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	return tags
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Age       int    `json:"age"`
	Timestamp string `json:"timestamp"`
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(username, "")
	if len(username) > 20 {
		username = username[:20]
	}
	return username
}

func TransformTimestamp(ts string) (string, error) {
	if ts == "" {
		return "", fmt.Errorf("timestamp cannot be empty")
	}
	return strings.ReplaceAll(ts, "/", "-"), nil
}

func ProcessUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if !ValidateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	data.Username = SanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age must be between 0 and 150")
	}

	transformedTS, err := TransformTimestamp(data.Timestamp)
	if err != nil {
		return nil, err
	}
	data.Timestamp = transformedTS

	return &data, nil
}

func main() {
	jsonData := `{"email":"test@example.com","username":"user_123!@#","age":25,"timestamp":"2023/12/31"}`
	processed, err := ProcessUserData([]byte(jsonData))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processed)
}