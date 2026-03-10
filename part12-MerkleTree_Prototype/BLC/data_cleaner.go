package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type Cleaner struct {
	trimSpaces bool
	lowercase  bool
}

func NewCleaner(trimSpaces, lowercase bool) *Cleaner {
	return &Cleaner{
		trimSpaces: trimSpaces,
		lowercase:  lowercase,
	}
}

func (c *Cleaner) ProcessString(input string) string {
	result := input
	if c.trimSpaces {
		result = strings.TrimSpace(result)
	}
	if c.lowercase {
		result = strings.ToLower(result)
	}
	return result
}

func CleanCSVFile(inputPath, outputPath string, cleaner *Cleaner) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedRecord[i] = cleaner.ProcessString(field)
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("error writing CSV: %w", err)
		}
	}

	return nil
}

func main() {
	cleaner := NewCleaner(true, true)
	err := CleanCSVFile("input.csv", "output.csv", cleaner)
	if err != nil {
		fmt.Printf("Error cleaning data: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Data cleaning completed successfully")
}