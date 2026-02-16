
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataCleaner struct {
	TrimSpaces bool
	RemoveEmpty bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		TrimSpaces: true,
		RemoveEmpty: true,
	}
}

func (dc *DataCleaner) CleanRow(row []string) []string {
	cleaned := make([]string, 0, len(row))
	
	for _, value := range row {
		processed := value
		
		if dc.TrimSpaces {
			processed = strings.TrimSpace(processed)
		}
		
		if dc.RemoveEmpty && processed == "" {
			continue
		}
		
		cleaned = append(cleaned, processed)
	}
	
	return cleaned
}

func (dc *DataCleaner) ProcessCSV(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()
	
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()
	
	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()
	
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record: %w", err)
		}
		
		cleanedRecord := dc.CleanRow(record)
		if len(cleanedRecord) == 0 {
			continue
		}
		
		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}
	
	return nil
}

func main() {
	cleaner := NewDataCleaner()
	
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}
	
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	
	err := cleaner.ProcessCSV(inputFile, outputFile)
	if err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Successfully cleaned data from %s to %s\n", inputFile, outputFile)
}