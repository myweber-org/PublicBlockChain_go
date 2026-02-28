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
}