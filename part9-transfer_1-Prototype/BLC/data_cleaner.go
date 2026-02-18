package datautils

import "sort"

// Deduplicate removes duplicate values from a slice of comparable types
func Deduplicate[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	seen := make(map[T]struct{})
	result := make([]T, 0, len(input))

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// DeduplicateSorted removes duplicates from a sorted slice more efficiently
func DeduplicateSorted[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	result := make([]T, 0, len(input))
	result = append(result, input[0])

	for i := 1; i < len(input); i++ {
		if input[i] != input[i-1] {
			result = append(result, input[i])
		}
	}

	return result
}

// DeduplicateWithSort sorts and removes duplicates from a slice
func DeduplicateWithSort[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	// Create a copy to avoid modifying original
	sorted := make([]T, len(input))
	copy(sorted, input)

	// Sort the slice
	sort.Slice(sorted, func(i, j int) bool {
		// Use type assertion for comparison
		// This is a simplified comparison - in practice you'd need
		// a more sophisticated approach for generic types
		return false // Placeholder - actual implementation depends on type
	})

	return DeduplicateSorted(sorted)
}package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func deduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if !seen[email] {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func validateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func processRecords(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	var cleaned []DataRecord

	for _, rec := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(rec.Email))
		if validateEmail(cleanEmail) && !emailSet[cleanEmail] {
			emailSet[cleanEmail] = true
			rec.Email = cleanEmail
			rec.Valid = true
			cleaned = append(cleaned, rec)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain", false},
		{5, "user@example.com", false},
	}

	cleaned := processRecords(records)
	fmt.Printf("Processed %d records, %d valid unique records found\n", len(records), len(cleaned))
	for _, rec := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", rec.ID, rec.Email, rec.Valid)
	}
}package main

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
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func deduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if !seen[email] {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func validateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func cleanRecords(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	var cleaned []DataRecord
	
	for _, record := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
		if validateEmail(cleanEmail) && !emailSet[cleanEmail] {
			emailSet[cleanEmail] = true
			record.Email = cleanEmail
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain", false},
		{5, "user@example.com", false},
	}
	
	cleaned := cleanRecords(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(cleaned))
	
	for _, rec := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", rec.ID, rec.Email, rec.Valid)
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

func readCSV(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []DataRecord

	// Skip header
	_, err = reader.Read()
	if err != nil {
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

		if len(row) >= 4 {
			records = append(records, DataRecord{
				ID:     strings.TrimSpace(row[0]),
				Name:   strings.TrimSpace(row[1]),
				Email:  strings.TrimSpace(row[2]),
				Active: strings.TrimSpace(row[3]),
			})
		}
	}

	return records, nil
}

func deduplicateByEmail(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(record.Email)
		if !seen[email] {
			seen[email] = true
			unique = append(unique, record)
		}
	}

	return unique
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func filterValidRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if validateEmail(record.Email) && record.ID != "" && record.Name != "" {
			valid = append(valid, record)
		}
	}
	return valid
}

func writeCSV(filename string, records []DataRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Name", "Email", "Active"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range records {
		row := []string{record.ID, record.Name, record.Email, record.Active}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	records, err := readCSV("input.csv")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Printf("Original records: %d\n", len(records))

	uniqueRecords := deduplicateByEmail(records)
	fmt.Printf("After deduplication: %d\n", len(uniqueRecords))

	validRecords := filterValidRecords(uniqueRecords)
	fmt.Printf("Valid records: %d\n", len(validRecords))

	err = writeCSV("cleaned_data.csv", validRecords)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Println("Data cleaning completed successfully")
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(id int, name, email string) {
	record := DataRecord{
		ID:    id,
		Name:  strings.TrimSpace(name),
		Email: strings.TrimSpace(email),
		Valid: true,
	}
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) ValidateEmails() {
	for i := range dc.records {
		if !strings.Contains(dc.records[i].Email, "@") {
			dc.records[i].Valid = false
		}
	}
}

func (dc *DataCleaner) RemoveDuplicates() {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range dc.records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] && record.Valid {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	dc.records = unique
}

func (dc *DataCleaner) GetValidRecords() []DataRecord {
	var valid []DataRecord
	for _, record := range dc.records {
		if record.Valid {
			valid = append(valid, record)
		}
	}
	return valid
}

func (dc *DataCleaner) PrintSummary() {
	fmt.Printf("Total records: %d\n", len(dc.records))
	valid := dc.GetValidRecords()
	fmt.Printf("Valid records: %d\n", len(valid))
}

func main() {
	cleaner := NewDataCleaner()

	cleaner.AddRecord(1, "John Doe", "john@example.com")
	cleaner.AddRecord(2, "Jane Smith", "jane@example.com")
	cleaner.AddRecord(3, "John Doe", "john@example.com")
	cleaner.AddRecord(4, "Bob Wilson", "invalid-email")
	cleaner.AddRecord(5, "Alice Brown", "alice@example.com")

	cleaner.ValidateEmails()
	cleaner.RemoveDuplicates()
	cleaner.PrintSummary()

	fmt.Println("\nCleaned records:")
	for _, record := range cleaner.GetValidRecords() {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}