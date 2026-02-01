package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	emailRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToLower(trimmed)
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, bool) {
	sanitizedName := dp.SanitizeInput(name)
	sanitizedEmail := dp.SanitizeInput(email)

	if sanitizedName == "" {
		return "Name cannot be empty", false
	}

	if !dp.ValidateEmail(sanitizedEmail) {
		return "Invalid email format", false
	}

	return "Data processed successfully", true
}package main

import (
	"errors"
	"strings"
)

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(username) > 20 {
		return errors.New("username must not exceed 20 characters")
	}
	if strings.Contains(username, " ") {
		return errors.New("username cannot contain spaces")
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func TransformToSlug(input string) string {
	slug := strings.ToLower(input)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	return slug
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

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
    file, err := os.Open(filePath)
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
            Valid: validateRecord(strings.TrimSpace(row[0]), strings.TrimSpace(row[2])),
        }

        records = append(records, record)
    }

    return records, nil
}

func validateRecord(id, email string) bool {
    if id == "" || email == "" {
        return false
    }
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FilterValidRecords(records []DataRecord) []DataRecord {
    var validRecords []DataRecord
    for _, record := range records {
        if record.Valid {
            validRecords = append(validRecords, record)
        }
    }
    return validRecords
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file_path>")
        return
    }

    records, err := ProcessCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        return
    }

    validRecords := FilterValidRecords(records)
    fmt.Printf("Total records: %d, Valid records: %d\n", len(records), len(validRecords))

    for _, record := range validRecords {
        fmt.Printf("ID: %s, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
    }
}
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value cannot be negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Value = record.Value * 1.1
	transformed.Tags = append(record.Tags, "processed")
	return transformed
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record))
	}
	return processed, nil
}

func GenerateSummary(records []DataRecord) string {
	if len(records) == 0 {
		return "No records to summarize"
	}
	
	var total float64
	var tagCount int
	for _, record := range records {
		total += record.Value
		tagCount += len(record.Tags)
	}
	
	avgValue := total / float64(len(records))
	avgTags := float64(tagCount) / float64(len(records))
	
	return fmt.Sprintf("Processed %d records. Average value: %.2f, Average tags per record: %.1f", 
		len(records), avgValue, avgTags)
}

func FilterByTag(records []DataRecord, tag string) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		for _, t := range record.Tags {
			if strings.EqualFold(t, tag) {
				filtered = append(filtered, record)
				break
			}
		}
	}
	return filtered
}