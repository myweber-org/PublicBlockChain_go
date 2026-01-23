package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	workers   int
	results   chan string
	errors    chan error
	waitGroup sync.WaitGroup
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		workers: workers,
		results: make(chan string, 100),
		errors:  make(chan error, 100),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		if lineCount%1000 == 0 {
			fp.results <- fmt.Sprintf("Processed %d lines from %s", lineCount, path)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	fp.results <- fmt.Sprintf("Completed processing %s: %d total lines", path, lineCount)
	return nil
}

func (fp *FileProcessor) Worker(fileQueue <-chan string) {
	defer fp.waitGroup.Done()
	for path := range fileQueue {
		if err := fp.ProcessFile(path); err != nil {
			fp.errors <- fmt.Errorf("error processing %s: %w", path, err)
		}
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string) ([]string, []error) {
	fileQueue := make(chan string, len(paths))

	for i := 0; i < fp.workers; i++ {
		fp.waitGroup.Add(1)
		go fp.Worker(fileQueue)
	}

	for _, path := range paths {
		fileQueue <- path
	}
	close(fileQueue)

	go func() {
		fp.waitGroup.Wait()
		close(fp.results)
		close(fp.errors)
	}()

	var allResults []string
	var allErrors []error

	for result := range fp.results {
		allResults = append(allResults, result)
	}

	for err := range fp.errors {
		allErrors = append(allErrors, err)
	}

	return allResults, allErrors
}

func FindFiles(rootDir, pattern string) ([]string, error) {
	var files []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return err
		}
		if matched {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: file_processor <directory> <pattern>")
		os.Exit(1)
	}

	rootDir := os.Args[1]
	pattern := os.Args[2]

	files, err := FindFiles(rootDir, pattern)
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No files found matching pattern")
		return
	}

	processor := NewFileProcessor(4)
	results, errors := processor.ProcessFiles(files)

	for _, result := range results {
		fmt.Println(result)
	}

	if len(errors) > 0 {
		fmt.Printf("\nEncountered %d errors:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	}

	fmt.Printf("\nProcessed %d files with %d errors\n", len(files), len(errors))
}