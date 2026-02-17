package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileStats struct {
	Path      string
	Size      int64
	LineCount int
	Processed bool
	Error     error
}

func processFile(path string, results chan<- FileStats, wg *sync.WaitGroup) {
	defer wg.Done()

	stats := FileStats{Path: path, Processed: false}

	file, err := os.Open(path)
	if err != nil {
		stats.Error = err
		results <- stats
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		stats.Error = err
		results <- stats
		return
	}
	stats.Size = info.Size()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		stats.Error = err
		results <- stats
		return
	}
	stats.LineCount = lineCount
	stats.Processed = true

	results <- stats
}

func collectFiles(dir string, extensions []string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		for _, ext := range extensions {
			if filepath.Ext(path) == ext {
				files = append(files, path)
				break
			}
		}
		return nil
	})
	return files, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	startTime := time.Now()
	dir := os.Args[1]
	extensions := []string{".txt", ".go", ".md", ".json"}

	files, err := collectFiles(dir, extensions)
	if err != nil {
		fmt.Printf("Error collecting files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No matching files found")
		return
	}

	results := make(chan FileStats, len(files))
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go processFile(file, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalSize := int64(0)
	totalLines := 0
	processedCount := 0
	errorCount := 0

	for stats := range results {
		if stats.Error != nil {
			fmt.Printf("Error processing %s: %v\n", stats.Path, stats.Error)
			errorCount++
			continue
		}
		if stats.Processed {
			fmt.Printf("Processed: %s | Size: %d bytes | Lines: %d\n",
				stats.Path, stats.Size, stats.LineCount)
			totalSize += stats.Size
			totalLines += stats.LineCount
			processedCount++
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\nSummary:\n")
	fmt.Printf("Files found: %d\n", len(files))
	fmt.Printf("Successfully processed: %d\n", processedCount)
	fmt.Printf("Errors: %d\n", errorCount)
	fmt.Printf("Total size: %d bytes\n", totalSize)
	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Processing time: %v\n", duration)
}package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
)

type FileProcessor struct {
	mu       sync.RWMutex
	filePath string
	lines    []string
}

func NewFileProcessor(path string) *FileProcessor {
	return &FileProcessor{
		filePath: path,
		lines:    make([]string, 0),
	}
}

func (fp *FileProcessor) Load() error {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	file, err := os.Open(fp.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fp.lines = make([]string, 0)

	for scanner.Scan() {
		fp.lines = append(fp.lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

func (fp *FileProcessor) ProcessLines(workerCount int) ([]string, error) {
	if len(fp.lines) == 0 {
		return nil, errors.New("no data loaded")
	}

	var wg sync.WaitGroup
	results := make([]string, len(fp.lines))
	chunkSize := (len(fp.lines) + workerCount - 1) / workerCount

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		start := i * chunkSize
		end := start + chunkSize
		if end > len(fp.lines) {
			end = len(fp.lines)
		}

		go func(workerID, s, e int) {
			defer wg.Done()
			for idx := s; idx < e; idx++ {
				processed := fmt.Sprintf("Worker-%d: %s", workerID, fp.lines[idx])
				results[idx] = processed
			}
		}(i, start, end)
	}

	wg.Wait()
	return results, nil
}

func (fp *FileProcessor) GetLineCount() int {
	fp.mu.RLock()
	defer fp.mu.RUnlock()
	return len(fp.lines)
}

func main() {
	processor := NewFileProcessor("sample.txt")
	
	if err := processor.Load(); err != nil {
		fmt.Printf("Error loading file: %v\n", err)
		return
	}

	fmt.Printf("Loaded %d lines\n", processor.GetLineCount())

	results, err := processor.ProcessLines(4)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	for _, line := range results {
		fmt.Println(line)
	}
}