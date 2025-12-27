
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileResult struct {
	Path     string
	Size     int64
	Lines    int
	Error    error
}

func processFile(path string, results chan<- FileResult, wg *sync.WaitGroup) {
	defer wg.Done()

	result := FileResult{Path: path}
	
	file, err := os.Open(path)
	if err != nil {
		result.Error = err
		results <- result
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		result.Error = err
		results <- result
		return
	}
	result.Size = info.Size()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	
	if err := scanner.Err(); err != nil {
		result.Error = err
		results <- result
		return
	}
	
	result.Lines = lineCount
	results <- result
}

func collectResults(results <-chan FileResult, totalFiles *int, totalSize *int64, totalLines *int) {
	for result := range results {
		if result.Error != nil {
			fmt.Printf("Error processing %s: %v\n", result.Path, result.Error)
			continue
		}
		
		*totalFiles++
		*totalSize += result.Size
		*totalLines += result.Lines
		
		fmt.Printf("Processed: %s (Size: %d bytes, Lines: %d)\n", 
			filepath.Base(result.Path), result.Size, result.Lines)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	root := os.Args[1]
	pattern := filepath.Join(root, "*.txt")
	
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No .txt files found in directory")
		return
	}

	fmt.Printf("Found %d .txt files to process\n\n", len(files))

	results := make(chan FileResult, len(files))
	var wg sync.WaitGroup

	startTime := time.Now()

	for _, file := range files {
		wg.Add(1)
		go processFile(file, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var totalFiles int
	var totalSize int64
	var totalLines int

	collectResults(results, &totalFiles, &totalSize, &totalLines)

	elapsed := time.Since(startTime)
	
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total files processed: %d\n", totalFiles)
	fmt.Printf("Total size: %d bytes\n", totalSize)
	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Processing time: %v\n", elapsed)
}