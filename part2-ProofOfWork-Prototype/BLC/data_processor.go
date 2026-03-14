
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	Email    string
	Username string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateProfile(p UserProfile) error {
	if !emailRegex.MatchString(p.Email) {
		return errors.New("invalid email format")
	}
	if strings.TrimSpace(p.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if p.Age < 0 || p.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(p *UserProfile) {
	p.Username = strings.ToLower(strings.TrimSpace(p.Username))
}

func ProcessUserProfile(p UserProfile) (UserProfile, error) {
	if err := ValidateProfile(p); err != nil {
		return p, err
	}
	TransformUsername(&p)
	return p, nil
}
package main

import "fmt"

func FilterAndDouble(nums []int, threshold int) []int {
    var result []int
    for _, num := range nums {
        if num > threshold {
            result = append(result, num*2)
        }
    }
    return result
}

func main() {
    input := []int{1, 5, 10, 15, 20}
    filtered := FilterAndDouble(input, 8)
    fmt.Println("Original:", input)
    fmt.Println("Filtered and doubled:", filtered)
}
package main

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	trimmed := strings.TrimSpace(input)

	pattern := `^[a-zA-Z0-9\s\-_\.@]+$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	if !re.MatchString(trimmed) {
		return "", nil
	}

	return trimmed, nil
}package main

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

func processCSV(filename string) ([]Record, error) {
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
            return nil, fmt.Errorf("invalid column count at line %d", lineNum)
        }

        id, err := strconv.Atoi(line[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
        }

        name := line[1]

        value, err := strconv.ParseFloat(line[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    return records, nil
}

func validateRecords(records []Record) error {
    seen := make(map[int]bool)
    for _, rec := range records {
        if rec.ID <= 0 {
            return fmt.Errorf("invalid ID %d", rec.ID)
        }
        if rec.Name == "" {
            return fmt.Errorf("empty name for ID %d", rec.ID)
        }
        if rec.Value < 0 {
            return fmt.Errorf("negative value for ID %d", rec.ID)
        }
        if seen[rec.ID] {
            return fmt.Errorf("duplicate ID %d", rec.ID)
        }
        seen[rec.ID] = true
    }
    return nil
}

func calculateStats(records []Record) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    for _, rec := range records {
        sum += rec.Value
    }
    average := sum / float64(len(records))

    var variance float64
    for _, rec := range records {
        diff := rec.Value - average
        variance += diff * diff
    }
    variance /= float64(len(records))

    return average, variance
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
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

		if len(line) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(line))
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := line[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
		lineNumber++
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) error {
	if len(records) == 0 {
		return fmt.Errorf("no records to validate")
	}

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
			return fmt.Errorf("negative value %f for record ID %d", record.Value, record.ID)
		}
	}

	return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, error) {
	if len(records) == 0 {
		return 0, 0, fmt.Errorf("cannot calculate statistics for empty dataset")
	}

	var sum, max float64
	for i, record := range records {
		sum += record.Value
		if i == 0 || record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max, nil
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
	return strings.TrimSpace(strings.ToLower(username))
}

func ProcessUserData(rawData []byte) (*UserProfile, error) {
	var profile UserProfile
	err := json.Unmarshal(rawData, &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if profile.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", profile.ID)
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
		"id": 1001,
		"username": "  JohnDoe  ",
		"email": "john@example.com",
		"age": 30,
		"active": true,
		"timestamp": "2024-01-15T10:30:00Z"
	}`

	profile, err := ProcessUserData([]byte(jsonData))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Printf("Processed user profile:\n")
	fmt.Printf("ID: %d\n", profile.ID)
	fmt.Printf("Username: %s\n", profile.Username)
	fmt.Printf("Email: %s\n", profile.Email)
	fmt.Printf("Age: %d\n", profile.Age)
	fmt.Printf("Active: %v\n", profile.Active)
	fmt.Printf("Timestamp: %s\n", profile.Timestamp)
}
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return normalized
}

func (dp *DataProcessor) ExtractAlphanumeric(input string) string {
	alphanumericRegex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return alphanumericRegex.ReplaceAllString(input, "")
}

func (dp *DataProcessor) Tokenize(input string) []string {
	cleaned := dp.CleanString(input)
	if cleaned == "" {
		return []string{}
	}
	return strings.Split(cleaned, " ")
}