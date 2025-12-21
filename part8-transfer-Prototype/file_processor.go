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

func (fp *FileProcessor) ProcessFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		processed := fmt.Sprintf("Line %d: %s", lineNumber, line)
		
		fp.mu.Lock()
		fp.results = append(fp.results, processed)
		fp.mu.Unlock()
		
		lineNumber++
	}

	return scanner.Err()
}

func (fp *FileProcessor) ProcessFilesConcurrently(filePaths []string) {
	for _, path := range filePaths {
		fp.wg.Add(1)
		go func(filePath string) {
			defer fp.wg.Done()
			if err := fp.ProcessFile(filePath); err != nil {
				fmt.Printf("Error processing %s: %v\n", filePath, err)
			}
		}(path)
	}
	fp.wg.Wait()
}

func (fp *FileProcessor) GetResults() []string {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.results
}

func main() {
	processor := NewFileProcessor()
	
	files := []string{"data1.txt", "data2.txt"}
	processor.ProcessFilesConcurrently(files)
	
	results := processor.GetResults()
	for _, result := range results {
		fmt.Println(result)
	}
}