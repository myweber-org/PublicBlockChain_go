package main

import (
	"errors"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(data UserData) UserData {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
	return data
}

func ProcessUserInput(rawUsername string, rawEmail string, rawAge int) (UserData, error) {
	userData := UserData{
		Username: rawUsername,
		Email:    rawEmail,
		Age:      rawAge,
	}

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	userData = TransformUsername(userData)
	return userData, nil
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateJSON(rawData []byte) (*UserData, error) {
	var user UserData
	err := json.Unmarshal(rawData, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if user.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.Name == "" {
		return nil, fmt.Errorf("user name cannot be empty")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user email cannot be empty")
	}

	return &user, nil
}

func main() {
	jsonInput := `{"id": 101, "name": "Alice", "email": "alice@example.com"}`
	user, err := ValidateJSON([]byte(jsonInput))
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}
	fmt.Printf("Validated user: %+v\n", user)
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

func normalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("invalid email format")
	}
	return email, nil
}

func validateUsername(username string) (string, error) {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return "", fmt.Errorf("username must be between 3 and 20 characters")
	}
	pattern := `^[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(pattern, username)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	return username, nil
}

func processUserData(rawData []byte) (*UserData, error) {
	var data UserData
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	email, err := normalizeEmail(data.Email)
	if err != nil {
		return nil, fmt.Errorf("email validation failed: %w", err)
	}
	data.Email = email

	username, err := validateUsername(data.Username)
	if err != nil {
		return nil, fmt.Errorf("username validation failed: %w", err)
	}
	data.Username = username

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age must be between 0 and 150")
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email": "  TEST@Example.COM ", "username": "user_123", "age": 25}`
	processedData, err := processUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processedData)
}
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataProcessor struct {
    InputPath  string
    OutputPath string
}

func NewDataProcessor(input, output string) *DataProcessor {
    return &DataProcessor{
        InputPath:  input,
        OutputPath: output,
    }
}

func (dp *DataProcessor) ValidateAndClean() error {
    inputFile, err := os.Open(dp.InputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(dp.OutputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    reader := csv.NewReader(inputFile)
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    header, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read header: %w", err)
    }

    if err := writer.Write(header); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    recordCount := 0
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            continue
        }

        if dp.isValidRecord(record) {
            cleanedRecord := dp.cleanRecord(record)
            if err := writer.Write(cleanedRecord); err != nil {
                return fmt.Errorf("failed to write record: %w", err)
            }
            recordCount++
        }
    }

    fmt.Printf("Processed %d valid records\n", recordCount)
    return nil
}

func (dp *DataProcessor) isValidRecord(record []string) bool {
    for _, field := range record {
        if strings.TrimSpace(field) == "" {
            return false
        }
    }
    return len(record) > 0
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
    cleaned := make([]string, len(record))
    for i, field := range record {
        cleaned[i] = strings.TrimSpace(field)
    }
    return cleaned
}

func main() {
    processor := NewDataProcessor("input.csv", "output.csv")
    if err := processor.ValidateAndClean(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}