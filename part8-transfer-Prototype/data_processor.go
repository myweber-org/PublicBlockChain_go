package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

func ParseAndValidateJSON(rawData []byte, requiredFields []string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var missingFields []string
	for _, field := range requiredFields {
		if _, exists := data[field]; !exists {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		return nil, ValidationError{
			Field:   "required_fields",
			Message: fmt.Sprintf("missing required fields: %s", strings.Join(missingFields, ", ")),
		}
	}

	for key, value := range data {
		if strVal, ok := value.(string); ok && strings.TrimSpace(strVal) == "" {
			return nil, ValidationError{
				Field:   key,
				Message: "field cannot be empty",
			}
		}
	}

	return data, nil
}

func main() {
	jsonData := []byte(`{"name": "test", "age": 25, "email": ""}`)
	required := []string{"name", "age", "email"}

	result, err := ParseAndValidateJSON(jsonData, required)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Validated data: %v\n", result)
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

		if !strings.Contains(record.Email, "@") {
			return nil, fmt.Errorf("invalid email format at line %d", lineNumber)
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	seenIDs := make(map[string]bool)

	for _, record := range records {
		if seenIDs[record.ID] {
			fmt.Printf("Duplicate ID found: %s\n", record.ID)
			continue
		}

		if record.Active != "true" && record.Active != "false" {
			fmt.Printf("Invalid active status for ID %s: %s\n", record.ID, record.Active)
			continue
		}

		seenIDs[record.ID] = true
		validRecords = append(validRecords, record)
	}

	return validRecords
}

func GenerateReport(records []DataRecord) {
	activeCount := 0
	for _, record := range records {
		if record.Active == "true" {
			activeCount++
		}
	}

	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Active records: %d\n", activeCount)
	fmt.Printf("Inactive records: %d\n", len(records)-activeCount)
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

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func sanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func transformUserData(rawData []byte) (*UserData, error) {
	var user UserData
	if err := json.Unmarshal(rawData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	if !validateEmail(user.Email) {
		return nil, fmt.Errorf("invalid email format: %s", user.Email)
	}

	user.Username = sanitizeUsername(user.Username)

	if user.Age < 0 || user.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", user.Age)
	}

	return &user, nil
}

func main() {
	rawJSON := `{"email":"test@example.com","username":"  john_doe  ","age":25}`
	user, err := transformUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}