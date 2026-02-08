
package main

import (
    "fmt"
    "strings"
)

type DataCleaner struct {
    duplicates map[string]bool
}

func NewDataCleaner() *DataCleaner {
    return &DataCleaner{
        duplicates: make(map[string]bool),
    }
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
    unique := []string{}
    for _, item := range items {
        normalized := strings.ToLower(strings.TrimSpace(item))
        if !dc.duplicates[normalized] {
            dc.duplicates[normalized] = true
            unique = append(unique, item)
        }
    }
    return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    return len(parts[0]) > 0 && len(parts[1]) > 0
}

func main() {
    cleaner := NewDataCleaner()
    
    data := []string{"apple", "banana", "Apple", "cherry", "banana", "  BANANA  "}
    unique := cleaner.RemoveDuplicates(data)
    fmt.Println("Unique items:", unique)
    
    emails := []string{"test@example.com", "invalid-email", "user@domain", "@nodomain", "domain@", ""}
    for _, email := range emails {
        fmt.Printf("Email '%s' valid: %v\n", email, cleaner.ValidateEmail(email))
    }
}
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Record struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func deduplicateRecords(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record
	for _, r := range records {
		key := fmt.Sprintf("%d|%s|%s", r.ID, strings.ToLower(r.Name), strings.ToLower(r.Email))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, r)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func cleanData(inputJSON string) (string, error) {
	var records []Record
	err := json.Unmarshal([]byte(inputJSON), &records)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	records = deduplicateRecords(records)
	var validRecords []Record
	for _, r := range records {
		if validateEmail(r.Email) {
			validRecords = append(validRecords, r)
		}
	}

	output, err := json.MarshalIndent(validRecords, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}
	return string(output), nil
}

func main() {
	input := `[
		{"id":1,"name":"John","email":"john@example.com"},
		{"id":2,"name":"Jane","email":"jane@example.com"},
		{"id":1,"name":"John","email":"john@example.com"},
		{"id":3,"name":"Bob","email":"invalid-email"}
	]`

	result, err := cleanData(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Cleaned data:")
	fmt.Println(result)
}