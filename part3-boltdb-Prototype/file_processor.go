package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	Workers   int
	BatchSize int
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		Workers:   workers,
		BatchSize: batchSize,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string, processFunc func(string) error) []error {
	var wg sync.WaitGroup
	errorChan := make(chan error, len(paths))
	pathChan := make(chan string, len(paths))

	for i := 0; i < fp.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathChan {
				if err := processFunc(path); err != nil {
					errorChan <- fmt.Errorf("processing %s: %w", path, err)
				}
			}
		}()
	}

	for _, path := range paths {
		pathChan <- path
	}
	close(pathChan)
	wg.Wait()
	close(errorChan)

	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}
	return errors
}

func readFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
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

func main() {
	processor := NewFileProcessor(4, 10)

	paths := []string{
		"data/file1.txt",
		"data/file2.txt",
		"data/file3.txt",
	}

	errors := processor.ProcessFiles(paths, func(path string) error {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		lines, err := readFileLines(absPath)
		if err != nil {
			return err
		}

		fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), len(lines))
		return nil
	})

	if len(errors) > 0 {
		fmt.Printf("Encountered %d errors:\n", len(errors))
		for _, err := range errors {
			fmt.Println(" -", err)
		}
	}
}
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu       sync.Mutex
	results  map[string]int
	wg       sync.WaitGroup
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		results: make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	fp.wg.Add(1)
	go func() {
		defer fp.wg.Done()
		
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", path, err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineCount := 0
		for scanner.Scan() {
			lineCount++
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error scanning file %s: %v\n", path, err)
			return
		}

		fp.mu.Lock()
		fp.results[path] = lineCount
		fp.mu.Unlock()
	}()

	return nil
}

func (fp *FileProcessor) Wait() {
	fp.wg.Wait()
}

func (fp *FileProcessor) GetResults() map[string]int {
	return fp.results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	processor := NewFileProcessor()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			if err := processor.ProcessFile(path); err != nil {
				fmt.Printf("Failed to process %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	processor.Wait()

	results := processor.GetResults()
	fmt.Println("File processing results:")
	for file, lines := range results {
		fmt.Printf("%s: %d lines\n", file, lines)
	}
}package main

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
		if info.IsDir() {
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

	if err != nil {
		return nil, err
	}
	return files, nil
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
		if stats.Processed {
			totalSize += stats.Size
			totalLines += stats.LineCount
			processedCount++
			fmt.Printf("Processed: %s (Size: %d bytes, Lines: %d)\n",
				stats.Path, stats.Size, stats.LineCount)
		} else {
			errorCount++
			fmt.Printf("Error processing %s: %v\n", stats.Path, stats.Error)
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\nSummary:\n")
	fmt.Printf("Total files processed: %d\n", processedCount)
	fmt.Printf("Files with errors: %d\n", errorCount)
	fmt.Printf("Total size: %d bytes\n", totalSize)
	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Processing time: %v\n", duration)
}