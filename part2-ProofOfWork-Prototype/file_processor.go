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
	workers int
	mu      sync.Mutex
	results map[string]int
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		workers: workers,
		results: make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFile(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanner error: %w", err)
	}

	fp.mu.Lock()
	fp.results[path] = lineCount
	fp.mu.Unlock()

	return lineCount, nil
}

func (fp *FileProcessor) ProcessDirectory(dir string) error {
	var wg sync.WaitGroup
	fileChan := make(chan string, fp.workers)
	errChan := make(chan error, fp.workers)

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileChan {
				_, err := fp.ProcessFile(path)
				if err != nil {
					errChan <- fmt.Errorf("processing %s: %w", path, err)
				}
			}
		}()
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		fileChan <- path
		return nil
	})

	close(fileChan)
	wg.Wait()
	close(errChan)

	if err != nil {
		return fmt.Errorf("walk error: %w", err)
	}

	var processingErrors []error
	for e := range errChan {
		processingErrors = append(processingErrors, e)
	}

	if len(processingErrors) > 0 {
		return fmt.Errorf("encountered %d errors during processing", len(processingErrors))
	}

	return nil
}

func (fp *FileProcessor) GetResults() map[string]int {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	
	resultCopy := make(map[string]int)
	for k, v := range fp.results {
		resultCopy[k] = v
	}
	return resultCopy
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	processor := NewFileProcessor(4)
	
	err := processor.ProcessDirectory(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing directory: %v\n", err)
		os.Exit(1)
	}

	results := processor.GetResults()
	fmt.Printf("Processed %d files\n", len(results))
	
	totalLines := 0
	for path, lines := range results {
		fmt.Printf("%s: %d lines\n", filepath.Base(path), lines)
		totalLines += lines
	}
	fmt.Printf("Total lines: %d\n", totalLines)
}