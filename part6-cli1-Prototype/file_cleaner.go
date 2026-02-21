package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const tempDir = "/tmp/myapp"
const retentionDays = 7

func main() {
	err := cleanOldFiles(tempDir, retentionDays)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string, days int) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	cutoffTime := time.Now().AddDate(0, 0, -days)
	removedCount := 0

	for _, file := range files {
		if file.ModTime().Before(cutoffTime) {
			fullPath := filepath.Join(dirPath, file.Name())
			err := os.Remove(fullPath)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", fullPath, err)
			} else {
				removedCount++
			}
		}
	}

	fmt.Printf("Removed %d files older than %d days\n", removedCount, days)
	return nil
}package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_cleaner <input_file> [output_file]")
		return
	}

	inputFile := os.Args[1]
	outputFile := "cleaned_" + inputFile
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}

	lines, err := readLines(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	uniqueLines := removeDuplicates(lines)

	err = writeLines(outputFile, uniqueLines)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Successfully cleaned file. Unique lines: %d -> %d\n", len(lines), len(uniqueLines))
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func removeDuplicates(lines []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}
	return result
}

func writeLines(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}