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
}package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	fp.mu.Lock()
	fp.results = append(fp.results, fmt.Sprintf("%s: %d lines", filepath.Base(path), lineCount))
	fp.mu.Unlock()

	return scanner.Err()
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fp.wg.Add(1)
		go func(filename string) {
			defer fp.wg.Done()
			fullPath := filepath.Join(dirPath, filename)
			if err := fp.ProcessFile(fullPath); err != nil {
				fmt.Printf("Error processing %s: %v\n", filename, err)
			}
		}(file.Name())
	}

	fp.wg.Wait()
	return nil
}

func (fp *FileProcessor) GetResults() []string {
	return fp.results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory_path>")
		os.Exit(1)
	}

	processor := NewFileProcessor()
	if err := processor.ProcessDirectory(os.Args[1]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processing results:")
	for _, result := range processor.GetResults() {
		fmt.Println(result)
	}
}package main

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