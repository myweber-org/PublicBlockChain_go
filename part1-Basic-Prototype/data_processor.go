
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
	ID    int
	Name  string
	Email string
	Score float64
}

func ProcessCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	lineNum := 0

	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNum, len(row))
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNum)
		}

		email := strings.TrimSpace(row[2])
		if !strings.Contains(email, "@") {
			return nil, fmt.Errorf("invalid email format at line %d", lineNum)
		}

		score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid score at line %d: %w", lineNum, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Email: email,
			Score: score,
		})
	}

	return records, nil
}

func CalculateAverage(records []Record) float64 {
	if len(records) == 0 {
		return 0.0
	}

	var total float64
	for _, record := range records {
		total += record.Score
	}
	return total / float64(len(records))
}

func FilterByScore(records []Record, minScore float64) []Record {
	var filtered []Record
	for _, record := range records {
		if record.Score >= minScore {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d records\n", len(records))
	fmt.Printf("Average score: %.2f\n", CalculateAverage(records))

	highScorers := FilterByScore(records, 80.0)
	fmt.Printf("Records with score >= 80: %d\n", len(highScorers))
}
package main

import (
    "errors"
    "fmt"
    "strings"
    "time"
)

type UserData struct {
    ID        int
    Name      string
    Email     string
    Age       int
    Timestamp time.Time
}

func ValidateUserData(data UserData) error {
    if data.ID <= 0 {
        return errors.New("invalid user ID")
    }
    if strings.TrimSpace(data.Name) == "" {
        return errors.New("name cannot be empty")
    }
    if !strings.Contains(data.Email, "@") {
        return errors.New("invalid email format")
    }
    if data.Age < 0 || data.Age > 150 {
        return errors.New("age must be between 0 and 150")
    }
    return nil
}

func TransformUserData(data UserData) UserData {
    return UserData{
        ID:        data.ID,
        Name:      strings.ToUpper(strings.TrimSpace(data.Name)),
        Email:     strings.ToLower(strings.TrimSpace(data.Email)),
        Age:       data.Age,
        Timestamp: time.Now(),
    }
}

func ProcessUserData(data UserData) (UserData, error) {
    if err := ValidateUserData(data); err != nil {
        return UserData{}, err
    }
    return TransformUserData(data), nil
}

func main() {
    sampleData := UserData{
        ID:    1001,
        Name:  "  john doe  ",
        Email: "JOHN@EXAMPLE.COM",
        Age:   30,
    }

    processedData, err := ProcessUserData(sampleData)
    if err != nil {
        fmt.Printf("Processing failed: %v\n", err)
        return
    }

    fmt.Printf("Original: %+v\n", sampleData)
    fmt.Printf("Processed: %+v\n", processedData)
}