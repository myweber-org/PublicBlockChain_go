
package main

import (
	"fmt"
	"math"
)

// FilterAndTransform processes a slice of integers, filters out values below threshold,
// and applies a transformation (square root of absolute value).
func FilterAndTransform(numbers []int, threshold int) []float64 {
	var result []float64
	for _, num := range numbers {
		if num > threshold {
			transformed := math.Sqrt(math.Abs(float64(num)))
			result = append(result, transformed)
		}
	}
	return result
}

func main() {
	input := []int{-10, 5, 3, 15, 8, -2, 25}
	threshold := 5
	output := FilterAndTransform(input, threshold)
	fmt.Printf("Processed slice: %v\n", output)
}
package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Comments string
}

func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(trimmed, "")
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ProcessUserData(data UserData) (UserData, error) {
	sanitizedData := UserData{
		Username: SanitizeInput(data.Username),
		Email:    SanitizeInput(data.Email),
		Comments: SanitizeInput(data.Comments),
	}

	if !ValidateEmail(sanitizedData.Email) {
		return sanitizedData, &InvalidEmailError{Email: sanitizedData.Email}
	}

	if len(sanitizedData.Username) < 3 {
		return sanitizedData, &InvalidUsernameError{Username: sanitizedData.Username}
	}

	return sanitizedData, nil
}

type InvalidEmailError struct {
	Email string
}

func (e *InvalidEmailError) Error() string {
	return "Invalid email format: " + e.Email
}

type InvalidUsernameError struct {
	Username string
}

func (e *InvalidUsernameError) Error() string {
	return "Username must be at least 3 characters long: " + e.Username
}
package main

import (
    "encoding/csv"
    "errors"
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
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
    lineNumber := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
        }

        lineNumber++
        if lineNumber == 1 {
            continue
        }

        if len(line) != 3 {
            return nil, errors.New("invalid column count")
        }

        id, err := strconv.Atoi(line[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID format: %w", err)
        }

        value, err := strconv.ParseFloat(line[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value format: %w", err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  line[1],
            Value: value,
        })
    }

    return records, nil
}

func ValidateRecords(records []Record) error {
    if len(records) == 0 {
        return errors.New("no records to validate")
    }

    seen := make(map[int]bool)
    for _, r := range records {
        if seen[r.ID] {
            return fmt.Errorf("duplicate ID found: %d", r.ID)
        }
        seen[r.ID] = true

        if r.Value < 0 {
            return fmt.Errorf("negative value for ID %d: %f", r.ID, r.Value)
        }
    }

    return nil
}

func CalculateStats(records []Record) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    var max float64 = records[0].Value

    for _, r := range records {
        sum += r.Value
        if r.Value > max {
            max = r.Value
        }
    }

    average := sum / float64(len(records))
    return average, max
}