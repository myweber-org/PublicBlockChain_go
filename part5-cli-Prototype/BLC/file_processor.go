
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileStats struct {
	Path     string
	Size     int64
	Lines    int
	Modified time.Time
}

func processFile(path string, wg *sync.WaitGroup, results chan<- FileStats) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", path, err)
		return
	}
	defer file.Close()

	stats := FileStats{Path: path}

	if fileInfo, err := file.Stat(); err == nil {
		stats.Size = fileInfo.Size()
		stats.Modified = fileInfo.ModTime()
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	stats.Lines = lineCount

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning %s: %v\n", path, err)
		return
	}

	results <- stats
}

func collectFiles(dir string, patterns []string) ([]string, error) {
	var files []string

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return nil, err
		}
		files = append(files, matches...)
	}

	return files, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory> [patterns...]")
		fmt.Println("Default patterns: *.txt, *.go, *.md")
		os.Exit(1)
	}

	dir := os.Args[1]
	patterns := []string{"*.txt", "*.go", "*.md"}
	if len(os.Args) > 2 {
		patterns = os.Args[2:]
	}

	files, err := collectFiles(dir, patterns)
	if err != nil {
		fmt.Printf("Error collecting files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No files found matching patterns")
		return
	}

	var wg sync.WaitGroup
	results := make(chan FileStats, len(files))

	fmt.Printf("Processing %d files...\n", len(files))

	for _, file := range files {
		wg.Add(1)
		go processFile(file, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalSize := int64(0)
	totalLines := 0
	fileCount := 0

	for stats := range results {
		fmt.Printf("%s: %d bytes, %d lines, modified %s\n",
			filepath.Base(stats.Path),
			stats.Size,
			stats.Lines,
			stats.Modified.Format("2006-01-02 15:04:05"))

		totalSize += stats.Size
		totalLines += stats.Lines
		fileCount++
	}

	fmt.Printf("\nSummary: Processed %d files, total size: %d bytes, total lines: %d\n",
		fileCount, totalSize, totalLines)
}