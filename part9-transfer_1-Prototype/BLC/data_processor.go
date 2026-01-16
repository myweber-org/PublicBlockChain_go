
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) bool {
	var js interface{}
	return json.Unmarshal(data, &js) == nil
}

// ParseJSONMap attempts to parse the byte slice into a map[string]interface{}.
func ParseJSONMap(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func main() {
	// Example usage
	validJSON := []byte(`{"name": "test", "value": 42}`)
	invalidJSON := []byte(`{name: test}`)

	fmt.Println("Valid JSON check:", ValidateJSON(validJSON))
	fmt.Println("Invalid JSON check:", ValidateJSON(invalidJSON))

	parsed, err := ParseJSONMap(validJSON)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed data: %v\n", parsed)
}
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

func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(trimmed, "")
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ProcessUserData(data UserData) (UserData, error) {
	data.Username = SanitizeInput(data.Username)
	data.Email = SanitizeInput(data.Email)
	data.Comments = SanitizeInput(data.Comments)

	if !ValidateEmail(data.Email) {
		return data, &InvalidEmailError{Email: data.Email}
	}

	if len(data.Username) < 3 || len(data.Username) > 50 {
		return data, &InvalidUsernameError{Username: data.Username}
	}

	return data, nil
}

type InvalidEmailError struct {
	Email string
}

func (e *InvalidEmailError) Error() string {
	return "Invalid email format: " + e.Email
}

type InvalidUsernameError struct {
	Username string
}

func (e *InvalidUsernameError) Error() string {
	return "Username must be between 3 and 50 characters: " + e.Username
}package main

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	pattern := `^[a-zA-Z0-9\s\.\-_@]+$`
	matched, err := regexp.MatchString(pattern, trimmed)
	if err != nil || !matched {
		return "", false
	}

	return trimmed, true
}
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

type UserData struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func ValidateAndParseJSON(input string) (*UserData, error) {
    trimmedInput := strings.TrimSpace(input)
    if trimmedInput == "" {
        return nil, fmt.Errorf("input is empty")
    }

    var data UserData
    err := json.Unmarshal([]byte(trimmedInput), &data)
    if err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }

    if data.Name == "" {
        return nil, fmt.Errorf("name field is required")
    }
    if data.Email == "" {
        return nil, fmt.Errorf("email field is required")
    }
    if data.Age <= 0 {
        return nil, fmt.Errorf("age must be a positive integer")
    }

    return &data, nil
}

func main() {
    jsonStr := `{"name": "John Doe", "email": "john@example.com", "age": 30}`
    user, err := ValidateAndParseJSON(jsonStr)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Parsed user: %+v\n", user)
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
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    if parts[0] == "" || parts[1] == "" {
        return false
    }
    return true
}

func FilterActiveUsers(records []DataRecord) []DataRecord {
    var activeUsers []DataRecord
    for _, record := range records {
        if strings.ToLower(record.Active) == "true" {
            activeUsers = append(activeUsers, record)
        }
    }
    return activeUsers
}

func GenerateReport(records []DataRecord) {
    fmt.Printf("Total records processed: %d\n", len(records))
    
    activeUsers := FilterActiveUsers(records)
    fmt.Printf("Active users: %d\n", len(activeUsers))

    invalidEmails := 0
    for _, record := range records {
        if !ValidateEmail(record.Email) {
            invalidEmails++
        }
    }
    fmt.Printf("Invalid email addresses: %d\n", invalidEmails)
}