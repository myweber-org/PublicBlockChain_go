
package main

import (
	"fmt"
)

// RemoveDuplicates removes duplicate strings from a slice while preserving order
func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_cleaner <input.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := strings.TrimSuffix(inputFile, ".csv") + "_cleaned.csv"

	err := removeDuplicates(inputFile, outputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Cleaned data saved to: %s\n", outputFile)
}

func removeDuplicates(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	seen := make(map[string]bool)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		key := strings.Join(record, "|")
		if !seen[key] {
			seen[key] = true
			err = writer.Write(record)
			if err != nil {
				return err
			}
		}
	}

	return nil
}