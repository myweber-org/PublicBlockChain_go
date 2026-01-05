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
	Age   int     `json:"age"`
	Score float64 `json:"score"`
}

func readCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	if _, err := reader.Read(); err != nil {
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

		age, _ := strconv.Atoi(row[1])
		score, _ := strconv.ParseFloat(row[2], 64)

		records = append(records, Record{
			Name:  row[0],
			Age:   age,
			Score: score,
		})
	}

	return records, nil
}

func convertToJSON(records []Record) (string, error) {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func processData(filename string) error {
	records, err := readCSV(filename)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("no records found")
	}

	jsonData, err := convertToJSON(records)
	if err != nil {
		return err
	}

	fmt.Println("Processed", len(records), "records")
	fmt.Println(jsonData)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	if err := processData(os.Args[1]); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
package main

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(trimmed, " ")
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func ContainsOnlyAlphanumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateAndTransform(data UserData) (UserData, error) {
	if strings.TrimSpace(data.Username) == "" {
		return data, fmt.Errorf("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return data, fmt.Errorf("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return data, fmt.Errorf("age must be between 0 and 150")
	}

	transformed := UserData{
		Username: strings.ToLower(strings.TrimSpace(data.Username)),
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Age:      data.Age,
	}
	return transformed, nil
}

func main() {
	testData := UserData{
		Username: "  TestUser  ",
		Email:    "TEST@EXAMPLE.COM",
		Age:      25,
	}

	result, err := ValidateAndTransform(testData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Original: %+v\n", testData)
	fmt.Printf("Processed: %+v\n", result)
}