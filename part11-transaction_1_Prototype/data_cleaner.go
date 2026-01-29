
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	input := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(input)
	fmt.Println("Original:", input)
	fmt.Println("Cleaned:", cleaned)
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
	Email string
	Phone string
}

type DataCleaner struct {
	records map[string]DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make(map[string]DataRecord),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) bool {
	key := record.ID + "|" + record.Email
	if _, exists := dc.records[key]; exists {
		return false
	}
	dc.records[key] = record
	return true
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (dc *DataCleaner) RemoveDuplicates() int {
	initialCount := len(dc.records)
	unique := make(map[string]DataRecord)
	
	for _, record := range dc.records {
		key := record.Email + "|" + record.Phone
		unique[key] = record
	}
	
	dc.records = unique
	return initialCount - len(dc.records)
}

func (dc *DataCleaner) LoadFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Email: strings.TrimSpace(row[1]),
			Phone: strings.TrimSpace(row[2]),
		}

		if dc.ValidateEmail(record.Email) {
			dc.AddRecord(record)
		}
	}
	return nil
}

func (dc *DataCleaner) ExportToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Email", "Phone"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range dc.records {
		row := []string{record.ID, record.Email, record.Phone}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func (dc *DataCleaner) Statistics() {
	fmt.Printf("Total records: %d\n", len(dc.records))
	
	emailDomains := make(map[string]int)
	for _, record := range dc.records {
		parts := strings.Split(record.Email, "@")
		if len(parts) == 2 {
			emailDomains[parts[1]]++
		}
	}
	
	fmt.Println("Email domain distribution:")
	for domain, count := range emailDomains {
		fmt.Printf("  %s: %d\n", domain, count)
	}
}

func main() {
	cleaner := NewDataCleaner()
	
	err := cleaner.LoadFromCSV("input.csv")
	if err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		return
	}
	
	removed := cleaner.RemoveDuplicates()
	fmt.Printf("Removed %d duplicate records\n", removed)
	
	cleaner.Statistics()
	
	err = cleaner.ExportToCSV("cleaned_data.csv")
	if err != nil {
		fmt.Printf("Error exporting data: %v\n", err)
		return
	}
	
	fmt.Println("Data cleaning completed successfully")
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSVData(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	cleanedHeaders := make([]string, len(headers))
	for i, h := range headers {
		cleanedHeaders[i] = strings.TrimSpace(strings.ToLower(h))
	}

	if err := writer.Write(cleanedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		cleanedRecord := make([]string, len(record))
		hasEmpty := false
		for i, field := range record {
			cleaned := strings.TrimSpace(field)
			if cleaned == "" {
				hasEmpty = true
			}
			cleanedRecord[i] = cleaned
		}

		if !hasEmpty {
			if err := writer.Write(cleanedRecord); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
			recordCount++
		}
	}

	fmt.Printf("Cleaned %d records, removed rows with empty fields\n", recordCount)
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSVData(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
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
    Email string
    Phone string
}

type DataCleaner struct {
    records []DataRecord
    seen    map[string]bool
}

func NewDataCleaner() *DataCleaner {
    return &DataCleaner{
        seen: make(map[string]bool),
    }
}

func (dc *DataCleaner) AddRecord(record DataRecord) bool {
    key := fmt.Sprintf("%s|%s", strings.ToLower(record.Email), record.Phone)
    if dc.seen[key] {
        return false
    }
    dc.seen[key] = true
    dc.records = append(dc.records, record)
    return true
}

func (dc *DataCleaner) ValidateRecords() []DataRecord {
    var valid []DataRecord
    for _, r := range dc.records {
        if strings.Contains(r.Email, "@") && len(r.Phone) >= 10 {
            valid = append(valid, r)
        }
    }
    return valid
}

func (dc *DataCleaner) GetUniqueCount() int {
    return len(dc.records)
}

func main() {
    cleaner := NewDataCleaner()
    
    records := []DataRecord{
        {1, "user@example.com", "1234567890"},
        {2, "user@example.com", "1234567890"},
        {3, "invalid-email", "9876543210"},
        {4, "another@test.com", "5551234567"},
    }
    
    for _, r := range records {
        added := cleaner.AddRecord(r)
        fmt.Printf("Record %d added: %v\n", r.ID, added)
    }
    
    fmt.Printf("Total unique records: %d\n", cleaner.GetUniqueCount())
    
    valid := cleaner.ValidateRecords()
    fmt.Printf("Valid records: %d\n", len(valid))
}
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}