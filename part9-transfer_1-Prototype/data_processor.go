
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

func (dp *DataProcessor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	cleaned := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return cleaned
}

func (dp *DataProcessor) NormalizeCase(input string, lower bool) string {
	if lower {
		return strings.ToLower(input)
	}
	return strings.ToUpper(input)
}

func (dp *DataProcessor) RemoveSpecialChars(input string, keepPattern string) string {
	if keepPattern == "" {
		keepPattern = `[^a-zA-Z0-9\s]`
	}
	re := regexp.MustCompile(keepPattern)
	return re.ReplaceAllString(input, "")
}

func main() {
	processor := NewDataProcessor()
	sample := "  Hello   World!  This  is  a  TEST.  "
	
	cleaned := processor.CleanInput(sample)
	normalized := processor.NormalizeCase(cleaned, true)
	final := processor.RemoveSpecialChars(normalized, `[^a-zA-Z\s]`)
	
	println("Original:", sample)
	println("Processed:", final)
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
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func sanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func processUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !validateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = sanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	jsonData := []byte(`{"email":"test@example.com","username":"  john_doe  ","age":25}`)
	processed, err := processUserData(jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processed)
}