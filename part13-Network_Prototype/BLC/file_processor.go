package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type FileProcessor struct {
	mu       sync.Mutex
	results  []string
	wg       sync.WaitGroup
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		results: make([]string, 0),
	}
}

func (fp *FileProcessor) ProcessFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		fp.wg.Add(1)
		go func(line string, num int) {
			defer fp.wg.Done()
			processed := fp.processLine(line, num)
			
			fp.mu.Lock()
			fp.results = append(fp.results, processed)
			fp.mu.Unlock()
		}(scanner.Text(), lineNum)
		lineNum++
	}

	fp.wg.Wait()
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	
	return nil
}

func (fp *FileProcessor) processLine(line string, lineNum int) string {
	return fmt.Sprintf("Line %d: %d characters processed", lineNum, len(line))
}

func (fp *FileProcessor) GetResults() []string {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <filename>")
		os.Exit(1)
	}

	processor := NewFileProcessor()
	
	err := processor.ProcessFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	results := processor.GetResults()
	for _, result := range results {
		fmt.Println(result)
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
	mu          sync.Mutex
	processed   int
	errors      []string
}

func (fp *FileProcessor) ProcessFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fp.mu.Lock()
		fp.errors = append(fp.errors, fmt.Sprintf("Failed to open %s: %v", path, err))
		fp.mu.Unlock()
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fp.mu.Lock()
		fp.errors = append(fp.errors, fmt.Sprintf("Error scanning %s: %v", path, err))
		fp.mu.Unlock()
		return
	}

	fp.mu.Lock()
	fp.processed++
	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
	fp.mu.Unlock()
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		errors: make([]string, 0),
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file1> <file2> ...")
		return
	}

	processor := NewFileProcessor()
	var wg sync.WaitGroup

	for _, filePath := range os.Args[1:] {
		wg.Add(1)
		go processor.ProcessFile(filePath, &wg)
	}

	wg.Wait()

	fmt.Printf("\nProcessing complete:\n")
	fmt.Printf("Files processed: %d\n", processor.processed)
	if len(processor.errors) > 0 {
		fmt.Printf("Errors encountered: %d\n", len(processor.errors))
		for _, err := range processor.errors {
			fmt.Printf("  - %s\n", err)
		}
	}
}