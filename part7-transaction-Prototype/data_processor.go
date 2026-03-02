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
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) != 3 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		name := row[1]
		if name == "" {
			continue
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func CalculateTotal(records []Record) float64 {
	var total float64
	for _, r := range records {
		total += r.Value
	}
	return total
}

func main() {
	records, err := ProcessCSV("data.csv")
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Total value: %.2f\n", CalculateTotal(records))
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
	fmt.Printf("Total valid records: %d\n", len(records))
	for i, record := range records {
		fmt.Printf("%d. ID: %s, Name: %s, Email: %s\n", i+1, record.ID, record.Name, record.Email)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file_path>")
		return
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	GenerateReport(records)
}
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Count int     `json:"count"`
}

func processCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 3 {
			continue
		}

		value, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			continue
		}

		count, err := strconv.Atoi(row[2])
		if err != nil {
			continue
		}

		record := Record{
			Name:  row[0],
			Value: value,
			Count: count,
		}
		records = append(records, record)
	}

	return records, nil
}

func convertToJSON(records []Record) (string, error) {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := processCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	jsonOutput, err := convertToJSON(records)
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(jsonOutput)
}
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        string
	Email     string
	Username  string
	Age       int
	Active    bool
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateProfile(p UserProfile) error {
	if p.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if !emailRegex.MatchString(p.Email) {
		return errors.New("invalid email format")
	}
	if len(p.Username) < 3 || len(p.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if p.Age < 0 || p.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformProfile(p UserProfile) UserProfile {
	p.Username = strings.ToLower(strings.TrimSpace(p.Username))
	p.Email = strings.ToLower(strings.TrimSpace(p.Email))
	return p
}

func ProcessUserProfile(p UserProfile) (UserProfile, error) {
	if err := ValidateProfile(p); err != nil {
		return UserProfile{}, err
	}
	return TransformProfile(p), nil
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID        int
	Name      string
	Value     float64
	Timestamp string
}

func parseCSVFile(filename string) ([]DataRecord, error) {
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

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		timestamp := strings.TrimSpace(row[3])
		if timestamp == "" {
			return nil, fmt.Errorf("empty timestamp at line %d", lineNumber)
		}

		record := DataRecord{
			ID:        id,
			Name:      name,
			Value:     value,
			Timestamp: timestamp,
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID %d found", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("negative value %f for ID %d", record.Value, record.ID)
		}
	}

	return nil
}

func calculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var maxValue float64
	count := len(records)

	for i, record := range records {
		sum += record.Value
		if i == 0 || record.Value > maxValue {
			maxValue = record.Value
		}
	}

	average := sum / float64(count)
	return average, maxValue, count
}

func processDataFile(filename string) error {
	records, err := parseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	if err := validateRecords(records); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	average, maxValue, count := calculateStatistics(records)

	fmt.Printf("Data processing complete\n")
	fmt.Printf("Records processed: %d\n", count)
	fmt.Printf("Average value: %.2f\n", average)
	fmt.Printf("Maximum value: %.2f\n", maxValue)

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	if err := processDataFile(filename); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
package data

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type Record struct {
	ID        string
	Email     string
	Timestamp time.Time
	Metadata  map[string]interface{}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeString(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func TransformRecord(record Record) (Record, error) {
	if err := ValidateEmail(record.Email); err != nil {
		return Record{}, err
	}

	normalizedEmail := NormalizeString(record.Email)
	processedMetadata := make(map[string]interface{})

	for key, value := range record.Metadata {
		switch v := value.(type) {
		case string:
			processedMetadata[key] = NormalizeString(v)
		default:
			processedMetadata[key] = v
		}
	}

	return Record{
		ID:        strings.ToUpper(record.ID),
		Email:     normalizedEmail,
		Timestamp: record.Timestamp.UTC(),
		Metadata:  processedMetadata,
	}, nil
}

func BatchProcess(records []Record) ([]Record, []error) {
	var processed []Record
	var errs []error

	for i, record := range records {
		transformed, err := TransformRecord(record)
		if err != nil {
			errs = append(errs, errors.New("record " + record.ID + " at index " + string(rune(i)) + ": " + err.Error()))
			continue
		}
		processed = append(processed, transformed)
	}

	return processed, errs
}
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Value = record.Value * 1.1
	transformed.Tags = append(record.Tags, "processed")
	transformed.Tags = normalizeTags(transformed.Tags)
	return transformed
}

func normalizeTags(tags []string) []string {
	uniqueTags := make(map[string]bool)
	var result []string
	
	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" && !uniqueTags[normalized] {
			uniqueTags[normalized] = true
			result = append(result, normalized)
		}
	}
	return result
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record))
	}
	
	return processed, nil
}

func main() {
	records := []DataRecord{
		{
			ID:        "rec001",
			Value:     100.0,
			Timestamp: time.Now(),
			Tags:      []string{"important", "TEST"},
		},
		{
			ID:        "rec002",
			Value:     200.5,
			Timestamp: time.Now().Add(-time.Hour),
			Tags:      []string{"urgent", "urgent"},
		},
	}
	
	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}
	
	fmt.Printf("Processed %d records successfully\n", len(processed))
	for _, rec := range processed {
		fmt.Printf("Record %s: Value=%.2f, Tags=%v\n", rec.ID, rec.Value, rec.Tags)
	}
}package main

import (
	"errors"
	"regexp"
	"strings"
)

// ValidateEmail checks if the provided string is a valid email address
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// NormalizePhone removes all non-digit characters from a phone number string
func NormalizePhone(phone string) (string, error) {
	re := regexp.MustCompile(`\D`)
	normalized := re.ReplaceAllString(phone, "")

	if len(normalized) < 10 {
		return "", errors.New("phone number too short")
	}

	return normalized, nil
}

// CapitalizeWords capitalizes the first letter of each word in a string
func CapitalizeWords(input string) string {
	words := strings.Fields(input)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

// GenerateSlug creates a URL-friendly slug from a title string
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}