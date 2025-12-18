package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	workers   int
	batchSize int
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string) (map[string]int, error) {
	if len(paths) == 0 {
		return nil, errors.New("no files to process")
	}

	var wg sync.WaitGroup
	results := make(map[string]int)
	resultChan := make(chan fileResult, len(paths))
	fileChan := make(chan string, len(paths))

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(fileChan, resultChan, &wg)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.err == nil {
			results[result.path] = result.lineCount
		}
	}

	return results, nil
}

type fileResult struct {
	path      string
	lineCount int
	err       error
}

func (fp *FileProcessor) worker(fileChan <-chan string, resultChan chan<- fileResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for path := range fileChan {
		lineCount, err := countLines(path)
		resultChan <- fileResult{
			path:      path,
			lineCount: lineCount,
			err:       err,
		}
	}
}

func countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	return lineCount, scanner.Err()
}

func main() {
	processor := NewFileProcessor(4, 10)

	files := []string{
		filepath.Join("testdata", "file1.txt"),
		filepath.Join("testdata", "file2.txt"),
	}

	start := time.Now()
	results, err := processor.ProcessFiles(files)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("Error processing files: %v\n", err)
		return
	}

	fmt.Printf("Processed %d files in %v\n", len(results), elapsed)
	for file, lines := range results {
		fmt.Printf("%s: %d lines\n", file, lines)
	}
}