
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