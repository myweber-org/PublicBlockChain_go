
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
}package main

import (
	"fmt"
	"strings"
)

type DataProcessor struct {
	filters []func(string) string
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		filters: []func(string) string{
			strings.TrimSpace,
			func(s string) string { return strings.ToLower(s) },
		},
	}
}

func (dp *DataProcessor) AddFilter(filter func(string) string) {
	dp.filters = append(dp.filters, filter)
}

func (dp *DataProcessor) Process(input string) string {
	result := input
	for _, filter := range dp.filters {
		result = filter(result)
	}
	return result
}

func main() {
	processor := NewDataProcessor()
	processor.AddFilter(func(s string) string {
		return strings.ReplaceAll(s, "test", "example")
	})

	input := "  TEST Data Processing Function  "
	output := processor.Process(input)
	fmt.Printf("Processed: '%s'\n", output)
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

func ProcessCSVFile(filename string) ([]DataRecord, error) {
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

		if len(row) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Name:  strings.TrimSpace(row[1]),
			Email: strings.TrimSpace(row[2]),
			Valid: validateRecord(row),
		}

		if record.Valid {
			records = append(records, record)
		}
	}

	return records, nil
}

func validateRecord(fields []string) bool {
	if len(fields) < 3 {
		return false
	}

	for _, field := range fields[:3] {
		if strings.TrimSpace(field) == "" {
			return false
		}
	}

	email := strings.TrimSpace(fields[2])
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func GenerateSummary(records []DataRecord) {
	validCount := 0
	for _, record := range records {
		if record.Valid {
			validCount++
		}
	}

	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Valid records: %d\n", validCount)
	fmt.Printf("Invalid records: %d\n", len(records)-validCount)
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

	GenerateSummary(records)

	for i, record := range records {
		if i < 5 && record.Valid {
			fmt.Printf("Sample record: %s - %s\n", record.ID, record.Name)
		}
	}
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
	ID      string
	Name    string
	Email   string
	Active  string
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
	headerSkipped := false

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		if !headerSkipped {
			headerSkipped = true
			continue
		}

		if len(row) < 4 {
			continue
		}

		record := DataRecord{
			ID:     strings.TrimSpace(row[0]),
			Name:   strings.TrimSpace(row[1]),
			Email:  strings.TrimSpace(row[2]),
			Active: strings.TrimSpace(row[3]),
		}

		if isValidRecord(record) {
			records = append(records, record)
		}
	}

	return records, nil
}

func isValidRecord(record DataRecord) bool {
	if record.ID == "" || record.Name == "" {
		return false
	}
	if !strings.Contains(record.Email, "@") {
		return false
	}
	return record.Active == "true" || record.Active == "false"
}

func FilterActiveRecords(records []DataRecord) []DataRecord {
	var active []DataRecord
	for _, record := range records {
		if record.Active == "true" {
			active = append(active, record)
		}
	}
	return active
}

func GenerateReport(records []DataRecord) {
	fmt.Printf("Total records processed: %d\n", len(records))
	active := FilterActiveRecords(records)
	fmt.Printf("Active records: %d\n", len(active))
	fmt.Printf("Inactive records: %d\n", len(records)-len(active))
}