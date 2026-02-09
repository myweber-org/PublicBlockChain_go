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
}package main

import (
	"errors"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(username) > 50 {
		return errors.New("username cannot exceed 50 characters")
	}
	return nil
}

func ValidateEmail(email string) error {
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func TransformUserData(username, email string, age int) (UserData, error) {
	if err := ValidateUsername(username); err != nil {
		return UserData{}, err
	}

	if err := ValidateEmail(email); err != nil {
		return UserData{}, err
	}

	normalizedEmail := NormalizeEmail(email)

	if age < 0 || age > 150 {
		return UserData{}, errors.New("age must be between 0 and 150")
	}

	return UserData{
		Username: username,
		Email:    normalizedEmail,
		Age:      age,
	}, nil
}package main

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
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
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

func parseCSV(filename string) ([]Record, error) {
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
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(line) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNum, len(line))
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  line[1],
			Value: value,
		})
		lineNum++
	}

	return records, nil
}

func validateRecords(records []Record) error {
	seenIDs := make(map[int]bool)
	for _, rec := range records {
		if rec.ID <= 0 {
			return fmt.Errorf("invalid record ID: %d must be positive", rec.ID)
		}
		if rec.Name == "" {
			return fmt.Errorf("record ID %d has empty name", rec.ID)
		}
		if rec.Value < 0 {
			return fmt.Errorf("record ID %d has negative value: %f", rec.ID, rec.Value)
		}
		if seenIDs[rec.ID] {
			return fmt.Errorf("duplicate record ID: %d", rec.ID)
		}
		seenIDs[rec.ID] = true
	}
	return nil
}

func calculateTotal(records []Record) float64 {
	total := 0.0
	for _, rec := range records {
		total += rec.Value
	}
	return total
}

func processDataFile(filename string) error {
	records, err := parseCSV(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	if err := validateRecords(records); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	total := calculateTotal(records)
	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Total value: %.2f\n", total)
	fmt.Printf("Average value: %.2f\n", total/float64(len(records)))

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	if err := processDataFile(os.Args[1]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
package data_processor

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func SanitizeUsername(input string) string {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	re := regexp.MustCompile(`[^a-z0-9_-]`)
	return re.ReplaceAllString(lower, "")
}

func ValidatePasswordStrength(password string) (bool, []string) {
	var issues []string
	
	if len(password) < 8 {
		issues = append(issues, "password must be at least 8 characters")
	}
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	
	if !hasUpper {
		issues = append(issues, "password must contain at least one uppercase letter")
	}
	if !hasLower {
		issues = append(issues, "password must contain at least one lowercase letter")
	}
	if !hasDigit {
		issues = append(issues, "password must contain at least one digit")
	}
	if !hasSpecial {
		issues = append(issues, "password must contain at least one special character")
	}
	
	return len(issues) == 0, issues
}

func NormalizePhoneNumber(phone string) string {
	re := regexp.MustCompile(`\D`)
	digits := re.ReplaceAllString(phone, "")
	
	if len(digits) == 10 {
		return digits
	}
	
	if len(digits) == 11 && digits[0] == '1' {
		return digits[1:]
	}
	
	return digits
}package data_processor

import (
	"regexp"
	"strings"
)

type Processor struct {
	whitespaceRegex *regexp.Regexp
}

func NewProcessor() *Processor {
	return &Processor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (p *Processor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	cleaned := p.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return cleaned
}

func (p *Processor) NormalizeCase(input string) string {
	return strings.ToLower(p.CleanString(input))
}

func (p *Processor) ExtractTokens(input string) []string {
	cleaned := p.CleanString(input)
	if cleaned == "" {
		return []string{}
	}
	return strings.Split(cleaned, " ")
}
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
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TrimAndLower(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func ContainsSQLInjection(input string) bool {
	keywords := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "UNION", "OR", "--"}
	upperInput := strings.ToUpper(input)
	for _, keyword := range keywords {
		if strings.Contains(upperInput, keyword) {
			return true
		}
	}
	return false
}