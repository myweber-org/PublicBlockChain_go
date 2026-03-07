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
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		fp.wg.Add(1)
		go func(ln int, text string) {
			defer fp.wg.Done()
			processed := fp.processLine(ln, text)
			fp.mu.Lock()
			fp.results = append(fp.results, processed)
			fp.mu.Unlock()
		}(lineNumber, line)
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	fp.wg.Wait()
	return nil
}

func (fp *FileProcessor) processLine(number int, text string) string {
	return fmt.Sprintf("Line %d: %d characters", number, len(text))
}

func (fp *FileProcessor) GetResults() []string {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return append([]string(nil), fp.results...)
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