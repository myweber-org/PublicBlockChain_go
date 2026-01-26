package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileStats struct {
	Path      string
	Size      int64
	LineCount int
	Processed bool
	Error     error
}

func processFile(path string, results chan<- FileStats, wg *sync.WaitGroup) {
	defer wg.Done()

	stats := FileStats{Path: path, Processed: false}

	file, err := os.Open(path)
	if err != nil {
		stats.Error = err
		results <- stats
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		stats.Error = err
		results <- stats
		return
	}
	stats.Size = info.Size()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		stats.Error = err
		results <- stats
		return
	}
	stats.LineCount = lineCount
	stats.Processed = true

	results <- stats
}

func collectFiles(dir string, extensions []string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		for _, ext := range extensions {
			if filepath.Ext(path) == ext {
				files = append(files, path)
				break
			}
		}
		return nil
	})
	return files, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	startTime := time.Now()
	dir := os.Args[1]
	extensions := []string{".txt", ".go", ".md", ".json"}

	files, err := collectFiles(dir, extensions)
	if err != nil {
		fmt.Printf("Error collecting files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No matching files found")
		return
	}

	results := make(chan FileStats, len(files))
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go processFile(file, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalSize := int64(0)
	totalLines := 0
	processedCount := 0
	errorCount := 0

	for stats := range results {
		if stats.Error != nil {
			fmt.Printf("Error processing %s: %v\n", stats.Path, stats.Error)
			errorCount++
			continue
		}
		if stats.Processed {
			fmt.Printf("Processed: %s | Size: %d bytes | Lines: %d\n",
				stats.Path, stats.Size, stats.LineCount)
			totalSize += stats.Size
			totalLines += stats.LineCount
			processedCount++
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\nSummary:\n")
	fmt.Printf("Files found: %d\n", len(files))
	fmt.Printf("Successfully processed: %d\n", processedCount)
	fmt.Printf("Errors: %d\n", errorCount)
	fmt.Printf("Total size: %d bytes\n", totalSize)
	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Processing time: %v\n", duration)
}