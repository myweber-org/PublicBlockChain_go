
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
    records := make([]Record, 0)

    for lineNum := 1; ; lineNum++ {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
        }

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d", lineNum)
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  row[1],
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found")
    }

    return records, nil
}

func ValidateRecords(records []Record) error {
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

func CalculateStats(records []Record) (float64, float64) {
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

func (dp *DataProcessor) SanitizeString(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToLower(trimmed)
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, string, bool) {
	sanitizedName := dp.SanitizeString(name)
	sanitizedEmail := dp.SanitizeString(email)
	isValid := dp.ValidateEmail(sanitizedEmail)
	return sanitizedName, sanitizedEmail, isValid
}