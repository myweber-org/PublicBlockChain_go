
package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ParseCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataRecord
	for i, row := range records {
		if len(row) != 3 {
			return nil, errors.New("invalid CSV format on line " + strconv.Itoa(i+1))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, errors.New("invalid ID on line " + strconv.Itoa(i+1))
		}

		name := row[1]
		if name == "" {
			return nil, errors.New("empty name on line " + strconv.Itoa(i+1))
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, errors.New("invalid value on line " + strconv.Itoa(i+1))
		}

		data = append(data, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return data, nil
}

func ValidateData(records []DataRecord) error {
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return errors.New("invalid ID: " + strconv.Itoa(record.ID))
		}
		if seenIDs[record.ID] {
			return errors.New("duplicate ID: " + strconv.Itoa(record.ID))
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return errors.New("negative value for ID: " + strconv.Itoa(record.ID))
		}
	}
	return nil
}

func ProcessCSVData(filePath string) ([]DataRecord, error) {
	data, err := ParseCSVFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := ValidateData(data); err != nil {
		return nil, err
	}

	return data, nil
}package data

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	alphaRegex = regexp.MustCompile(`^[a-zA-Z\s]+$`)
)

func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidateAlpha(input string) bool {
	return alphaRegex.MatchString(input)
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
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

    activeUsers := FilterActiveUsers(records)
    fmt.Printf("Total records: %d\n", len(records))
    fmt.Printf("Active users: %d\n", len(activeUsers))

    for _, user := range activeUsers {
        fmt.Printf("ID: %s, Name: %s\n", user.ID, user.Name)
    }
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func ProcessCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		lineNum++
		if lineNum == 1 {
			continue
		}

		if len(line) != 3 {
			return nil, fmt.Errorf("invalid column count on line %d", lineNum)
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID on line %d: %w", lineNum, err)
		}

		name := line[1]
		if name == "" {
			return nil, fmt.Errorf("empty name on line %d", lineNum)
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value on line %d: %w", lineNum, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func ValidateRecords(records []Record) error {
	seenIDs := make(map[int]bool)
	for _, rec := range records {
		if rec.ID <= 0 {
			return fmt.Errorf("invalid ID %d", rec.ID)
		}
		if seenIDs[rec.ID] {
			return fmt.Errorf("duplicate ID %d", rec.ID)
		}
		seenIDs[rec.ID] = true
	}
	return nil
}

func CalculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, rec := range records {
		sum += rec.Value
		if rec.Value > max {
			max = rec.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
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

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserData(rawData []byte) (*UserProfile, error) {
	var profile UserProfile
	err := json.Unmarshal(rawData, &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	if profile.Username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	profile.Username = TransformUsername(profile.Username)

	if !ValidateEmail(profile.Email) {
		return nil, fmt.Errorf("invalid email format: %s", profile.Email)
	}

	if profile.Age < 0 || profile.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", profile.Age)
	}

	return &profile, nil
}

func main() {
	jsonData := `{
		"id": 42,
		"username": "  JohnDoe  ",
		"email": "john@example.com",
		"age": 30,
		"active": true,
		"timestamp": "2024-01-15T10:30:00Z"
	}`

	processedProfile, err := ProcessUserData([]byte(jsonData))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Printf("Processed User Profile:\n")
	fmt.Printf("ID: %d\n", processedProfile.ID)
	fmt.Printf("Username: %s\n", processedProfile.Username)
	fmt.Printf("Email: %s\n", processedProfile.Email)
	fmt.Printf("Age: %d\n", processedProfile.Age)
	fmt.Printf("Active: %v\n", processedProfile.Active)
	fmt.Printf("Timestamp: %s\n", processedProfile.Timestamp)
}
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateJSON(data []byte) (*User, error) {
	var user User
	err := json.Unmarshal(data, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if user.Name == "" {
		return nil, fmt.Errorf("name field is required")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("email field is required")
	}
	if user.ID <= 0 {
		return nil, fmt.Errorf("id must be a positive integer")
	}

	return &user, nil
}

func main() {
	jsonData := []byte(`{"id": 123, "name": "John Doe", "email": "john@example.com"}`)

	user, err := ValidateJSON(jsonData)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	fmt.Printf("Valid user: %+v\n", user)
}