
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
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

func (fp *FileProcessor) ProcessFiles(paths []string, processor func(string) error) []error {
	var wg sync.WaitGroup
	errorChan := make(chan error, len(paths))
	pathChan := make(chan string, len(paths))

	for i := 0; i < fp.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathChan {
				if err := processor(path); err != nil {
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

func (fp *FileProcessor) ReadLines(path string) ([]string, error) {
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

func (fp *FileProcessor) BatchProcess(dir string, pattern string) error {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return err
	}

	batches := fp.createBatches(matches)
	for _, batch := range batches {
		errors := fp.ProcessFiles(batch, func(path string) error {
			content, err := fp.ReadLines(path)
			if err != nil {
				return err
			}
			fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), len(content))
			return nil
		})

		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Printf("Error: %v\n", err)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func (fp *FileProcessor) createBatches(items []string) [][]string {
	var batches [][]string
	for i := 0; i < len(items); i += fp.BatchSize {
		end := i + fp.BatchSize
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}
	return batches
}

func main() {
	processor := NewFileProcessor(4, 10)
	if err := processor.BatchProcess(".", "*.txt"); err != nil {
		fmt.Printf("Batch processing failed: %v\n", err)
	}
}