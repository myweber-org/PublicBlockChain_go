
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
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct{}

func (dc DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range items {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (dc DataCleaner) TrimWhitespace(items []string) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func main() {
	cleaner := DataCleaner{}
	data := []string{" apple ", "banana", " apple ", "  cherry  ", "banana"}

	fmt.Println("Original:", data)
	trimmed := cleaner.TrimWhitespace(data)
	fmt.Println("Trimmed:", trimmed)
	unique := cleaner.RemoveDuplicates(trimmed)
	fmt.Println("Unique:", unique)
}package main

import (
	"fmt"
	"strings"
)

func DeduplicateStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func NormalizeWhitespace(input string) string {
	words := strings.Fields(input)
	return strings.Join(words, " ")
}

func main() {
	data := []string{"apple", "banana", "apple", "cherry", "banana"}
	unique := DeduplicateStrings(data)
	fmt.Println("Deduplicated:", unique)

	text := "   Hello    world!   This   is   a   test.   "
	normalized := NormalizeWhitespace(text)
	fmt.Println("Normalized:", normalized)
}
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

import (
    "fmt"
    "strings"
)

type DataCleaner struct {
    duplicatesRemoved int
}

func NewDataCleaner() *DataCleaner {
    return &DataCleaner{duplicatesRemoved: 0}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
    seen := make(map[string]bool)
    unique := []string{}
    
    for _, item := range items {
        if !seen[item] {
            seen[item] = true
            unique = append(unique, item)
        } else {
            dc.duplicatesRemoved++
        }
    }
    return unique
}

func (dc *DataCleaner) NormalizeText(items []string) []string {
    normalized := make([]string, len(items))
    for i, item := range items {
        normalized[i] = strings.ToLower(strings.TrimSpace(item))
    }
    return normalized
}

func (dc *DataCleaner) GetStats() string {
    return fmt.Sprintf("Duplicates removed: %d", dc.duplicatesRemoved)
}

func main() {
    cleaner := NewDataCleaner()
    
    data := []string{"Apple", "apple", "Banana", "  Banana  ", "Apple", "Cherry"}
    
    normalized := cleaner.NormalizeText(data)
    unique := cleaner.RemoveDuplicates(normalized)
    
    fmt.Println("Original data:", data)
    fmt.Println("Cleaned data:", unique)
    fmt.Println(cleaner.GetStats())
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct{}

func (dc DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range items {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (dc DataCleaner) TrimWhitespace(items []string) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func main() {
	cleaner := DataCleaner{}
	data := []string{"  apple ", "banana", "  apple ", " cherry", "banana "}

	fmt.Println("Original:", data)
	trimmed := cleaner.TrimWhitespace(data)
	fmt.Println("Trimmed:", trimmed)
	unique := cleaner.RemoveDuplicates(trimmed)
	fmt.Println("Cleaned:", unique)
}
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
	input := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(input)
	fmt.Printf("Original: %v\n", input)
	fmt.Printf("Cleaned: %v\n", cleaned)
}
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataCleaner struct {
    inputPath  string
    outputPath string
    seenRows   map[string]bool
}

func NewDataCleaner(input, output string) *DataCleaner {
    return &DataCleaner{
        inputPath:  input,
        outputPath: output,
        seenRows:   make(map[string]bool),
    }
}

func (dc *DataCleaner) RemoveDuplicates() error {
    inputFile, err := os.Open(dc.inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(dc.outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    reader := csv.NewReader(inputFile)
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    header, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read header: %w", err)
    }

    if err := writer.Write(header); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read record: %w", err)
        }

        rowKey := strings.Join(record, "|")
        if !dc.seenRows[rowKey] {
            dc.seenRows[rowKey] = true
            if err := writer.Write(record); err != nil {
                return fmt.Errorf("failed to write record: %w", err)
            }
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    cleaner := NewDataCleaner(os.Args[1], os.Args[2])
    if err := cleaner.RemoveDuplicates(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
}