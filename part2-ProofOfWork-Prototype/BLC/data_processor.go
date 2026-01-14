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

func ProcessCSVFile(filename string) ([]Record, error) {
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
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNum)
		}

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

func CalculateStats(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var max float64 = records[0].Value
	count := len(records)

	for _, r := range records {
		sum += r.Value
		if r.Value > max {
			max = r.Value
		}
	}

	average := sum / float64(count)
	return average, max, count
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

	avg, max, count := CalculateStats(records)
	fmt.Printf("Processed %d records\n", count)
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)
}package main

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

func ParseUserJSON(data []byte) (*User, error) {
	var user User
	err := json.Unmarshal(data, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user JSON: %w", err)
	}
	return &user, nil
}

func ValidateUser(user *User) error {
	if user.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.Name == "" {
		return fmt.Errorf("user name cannot be empty")
	}
	if user.Email == "" {
		return fmt.Errorf("user email cannot be empty")
	}
	return nil
}

func ProcessUserData(jsonData []byte) (*User, error) {
	user, err := ParseUserJSON(jsonData)
	if err != nil {
		return nil, err
	}

	err = ValidateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func main() {
	jsonStr := `{"id": 123, "name": "John Doe", "email": "john@example.com"}`
	user, err := ProcessUserData([]byte(jsonStr))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Processed user: ID=%d, Name=%s, Email=%s\n", user.ID, user.Name, user.Email)
}package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	emailRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return &DataProcessor{emailRegex: regex}
}

func (dp *DataProcessor) SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	return input
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, bool) {
	sanitizedName := dp.SanitizeInput(name)
	if !dp.ValidateEmail(email) {
		return "", false
	}
	return sanitizedName, true
}