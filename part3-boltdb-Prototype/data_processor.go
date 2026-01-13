package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitized := re.ReplaceAllString(input, "")
	return strings.TrimSpace(sanitized)
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

func TrimAndLower(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func ContainsSQLInjection(input string) bool {
	patterns := []string{
		`(?i)select.*from`,
		`(?i)insert.*into`,
		`(?i)update.*set`,
		`(?i)delete.*from`,
		`(?i)drop.*table`,
		`(?i)union.*select`,
		`--`,
		`;`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(input) {
			return true
		}
	}
	return false
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

func ValidateAndParseJSON(rawData []byte) (*UserData, error) {
	var user UserData
	if err := json.Unmarshal(rawData, &user); err != nil {
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
	jsonStr := `{"id": 123, "name": "John Doe", "email": "john@example.com"}`
	user, err := ValidateAndParseJSON([]byte(jsonStr))
	if err != nil {
		log.Fatalf("Error: %v", err)
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
	ID    string
	Name  string
	Email string
	Valid bool
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
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
			Email: strings.TrimSpace(row[2]),
			Valid: validateRecord(strings.TrimSpace(row[0]), strings.TrimSpace(row[2])),
		}

		if record.Valid {
			records = append(records, record)
		}
	}

	return records, nil
}

func validateRecord(id, email string) bool {
	if id == "" || email == "" {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func GenerateReport(records []DataRecord) {
	fmt.Printf("Data Processing Report\n")
	fmt.Printf("======================\n")
	fmt.Printf("Total valid records: %d\n", len(records))

	emailDomains := make(map[string]int)
	for _, record := range records {
		parts := strings.Split(record.Email, "@")
		if len(parts) == 2 {
			emailDomains[parts[1]]++
		}
	}

	fmt.Printf("\nEmail domain distribution:\n")
	for domain, count := range emailDomains {
		fmt.Printf("  %s: %d\n", domain, count)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := ProcessCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	GenerateReport(records)
}