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