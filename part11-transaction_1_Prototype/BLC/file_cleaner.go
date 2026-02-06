
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	tempDir := os.TempDir()
	cutoffTime := time.Now().AddDate(0, 0, -7)
	var removedCount int

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			if err := os.RemoveAll(path); err == nil {
				removedCount++
				fmt.Printf("Removed: %s\n", path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	fmt.Printf("Cleaning completed. Removed %d items.\n", removedCount)
}package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_cleaner.go <input_file>")
		return
	}

	inputFile := os.Args[1]
	outputFile := "cleaned_" + inputFile

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	seen := make(map[string]bool)
	var uniqueLines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] {
			seen[line] = true
			uniqueLines = append(uniqueLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	for _, line := range uniqueLines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()

	fmt.Printf("Duplicate lines removed. Cleaned file saved as: %s\n", outputFile)
}