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
	mu       sync.Mutex
	fileList []string
	errors   []error
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		fileList: make([]string, 0),
		errors:   make([]error, 0),
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			fp.mu.Lock()
			fp.fileList = append(fp.fileList, path)
			fp.mu.Unlock()
		}
		return nil
	})
	return err
}

func (fp *FileProcessor) CountLines() map[string]int {
	var wg sync.WaitGroup
	result := make(map[string]int)
	resultMu := sync.Mutex{}

	for _, file := range fp.fileList {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			lines, err := countFileLines(f)
			if err != nil {
				fp.mu.Lock()
				fp.errors = append(fp.errors, fmt.Errorf("file %s: %w", f, err))
				fp.mu.Unlock()
				return
			}
			resultMu.Lock()
			result[f] = lines
			resultMu.Unlock()
		}(file)
	}
	wg.Wait()
	return result
}

func countFileLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lineCount := 0
	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return lineCount, err
		}
		lineCount++
	}
	return lineCount, nil
}

func (fp *FileProcessor) GetErrors() []error {
	return fp.errors
}

func main() {
	processor := NewFileProcessor()
	
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}
	
	dir := os.Args[1]
	err := processor.ProcessDirectory(dir)
	if err != nil {
		fmt.Printf("Error processing directory: %v\n", err)
		os.Exit(1)
	}
	
	results := processor.CountLines()
	fmt.Println("Line counts:")
	for file, count := range results {
		fmt.Printf("%s: %d lines\n", file, count)
	}
	
	if errors := processor.GetErrors(); len(errors) > 0 {
		fmt.Println("\nErrors encountered:")
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	}
}