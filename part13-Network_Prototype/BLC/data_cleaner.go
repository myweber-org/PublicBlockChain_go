package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6, 7}
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
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

type Record struct {
    ID      int
    Name    string
    Email   string
    Active  bool
    Score   float64
}

func cleanEmail(email string) string {
    return strings.ToLower(strings.TrimSpace(email))
}

func validateRecord(record []string) (*Record, error) {
    if len(record) != 5 {
        return nil, fmt.Errorf("invalid record length: %d", len(record))
    }

    id, err := strconv.Atoi(strings.TrimSpace(record[0]))
    if err != nil {
        return nil, fmt.Errorf("invalid ID: %v", err)
    }

    name := strings.TrimSpace(record[1])
    if name == "" {
        return nil, fmt.Errorf("name cannot be empty")
    }

    email := cleanEmail(record[2])
    if !strings.Contains(email, "@") {
        return nil, fmt.Errorf("invalid email format")
    }

    active, err := strconv.ParseBool(strings.TrimSpace(record[3]))
    if err != nil {
        return nil, fmt.Errorf("invalid active status: %v", err)
    }

    score, err := strconv.ParseFloat(strings.TrimSpace(record[4]), 64)
    if err != nil {
        return nil, fmt.Errorf("invalid score: %v", err)
    }

    return &Record{
        ID:     id,
        Name:   name,
        Email:  email,
        Active: active,
        Score:  score,
    }, nil
}

func processCSVFile(filename string) ([]Record, []error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, []error{err}
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    var records []Record
    var errors []error
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: read error: %v", lineNumber, err))
            continue
        }

        if lineNumber == 1 {
            continue
        }

        record, err := validateRecord(row)
        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: %v", lineNumber, err))
            continue
        }

        records = append(records, *record)
    }

    return records, errors
}

func generateSummary(records []Record) {
    var totalScore float64
    activeCount := 0
    emailDomains := make(map[string]int)

    for _, record := range records {
        totalScore += record.Score
        if record.Active {
            activeCount++
        }

        parts := strings.Split(record.Email, "@")
        if len(parts) == 2 {
            emailDomains[parts[1]]++
        }
    }

    avgScore := totalScore / float64(len(records))
    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Active records: %d\n", activeCount)
    fmt.Printf("Average score: %.2f\n", avgScore)
    fmt.Println("Email domain distribution:")
    for domain, count := range emailDomains {
        fmt.Printf("  %s: %d\n", domain, count)
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run data_cleaner.go <csv_file>")
        os.Exit(1)
    }

    filename := os.Args[1]
    records, errors := processCSVFile(filename)

    if len(errors) > 0 {
        fmt.Printf("Encountered %d errors during processing:\n", len(errors))
        for _, err := range errors {
            fmt.Printf("  - %v\n", err)
        }
    }

    if len(records) > 0 {
        fmt.Println("\nSuccessfully processed records:")
        generateSummary(records)

        fmt.Println("\nFirst 5 valid records:")
        for i := 0; i < len(records) && i < 5; i++ {
            fmt.Printf("  ID: %d, Name: %s, Email: %s, Active: %v, Score: %.1f\n",
                records[i].ID, records[i].Name, records[i].Email,
                records[i].Active, records[i].Score)
        }
    }
}