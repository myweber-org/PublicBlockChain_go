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

type FileProcessor struct {
	workers   int
	batchSize int
	mu        sync.Mutex
	stats     ProcessingStats
}

type ProcessingStats struct {
	FilesProcessed int
	BytesRead      int64
	Errors         int
	StartTime      time.Time
	EndTime        time.Time
}

type FileTask struct {
	Path    string
	Content []byte
	Err     error
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	if workers < 1 {
		workers = 4
	}
	if batchSize < 1 {
		batchSize = 10
	}
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
		stats:     ProcessingStats{},
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	fp.stats.StartTime = time.Now()
	defer func() {
		fp.stats.EndTime = time.Now()
	}()

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	taskChan := make(chan string, len(files))
	resultChan := make(chan FileTask, len(files))
	var wg sync.WaitGroup

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(taskChan, resultChan, &wg)
	}

	for _, file := range files {
		if !file.IsDir() {
			fullPath := filepath.Join(dirPath, file.Name())
			taskChan <- fullPath
		}
	}
	close(taskChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		fp.mu.Lock()
		fp.stats.FilesProcessed++
		if result.Err != nil {
			fp.stats.Errors++
			fmt.Printf("Error processing %s: %v\n", result.Path, result.Err)
		} else {
			fp.stats.BytesRead += int64(len(result.Content))
		}
		fp.mu.Unlock()
	}

	return nil
}

func (fp *FileProcessor) worker(taskChan <-chan string, resultChan chan<- FileTask, wg *sync.WaitGroup) {
	defer wg.Done()

	for path := range taskChan {
		content, err := fp.readFileInBatches(path)
		resultChan <- FileTask{
			Path:    path,
			Content: content,
			Err:     err,
		}
	}
}

func (fp *FileProcessor) readFileInBatches(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content []byte
	reader := bufio.NewReader(file)
	buffer := make([]byte, fp.batchSize*1024)

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			content = append(content, buffer[:n]...)
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return content, err
		}
	}

	return content, nil
}

func (fp *FileProcessor) GetStats() ProcessingStats {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.stats
}

func (fp *FileProcessor) PrintSummary() {
	stats := fp.GetStats()
	duration := stats.EndTime.Sub(stats.StartTime)

	fmt.Println("\n=== Processing Summary ===")
	fmt.Printf("Files processed: %d\n", stats.FilesProcessed)
	fmt.Printf("Total bytes read: %d\n", stats.BytesRead)
	fmt.Printf("Errors encountered: %d\n", stats.Errors)
	fmt.Printf("Processing time: %v\n", duration.Round(time.Millisecond))
	if stats.FilesProcessed > 0 && duration > 0 {
		throughput := float64(stats.BytesRead) / duration.Seconds()
		fmt.Printf("Throughput: %.2f bytes/second\n", throughput)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory_path>")
		fmt.Println("Example: file_processor ./data")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	processor := NewFileProcessor(4, 10)

	fmt.Printf("Processing files in: %s\n", dirPath)
	err := processor.ProcessDirectory(dirPath)
	if err != nil {
		fmt.Printf("Failed to process directory: %v\n", err)
		os.Exit(1)
	}

	processor.PrintSummary()
}package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileResult struct {
	Path     string
	Size     int64
	Lines    int
	Error    error
}

func processFile(path string, results chan<- FileResult, wg *sync.WaitGroup) {
	defer wg.Done()

	result := FileResult{Path: path}
	file, err := os.Open(path)
	if err != nil {
		result.Error = err
		results <- result
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		result.Error = err
		results <- result
		return
	}
	result.Size = info.Size()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		result.Error = err
		results <- result
		return
	}
	result.Lines = lineCount

	results <- result
}

func collectResults(results <-chan FileResult, totalFiles *int, totalSize *int64, totalLines *int) {
	for result := range results {
		if result.Error != nil {
			fmt.Printf("Error processing %s: %v\n", result.Path, result.Error)
			continue
		}
		*totalFiles++
		*totalSize += result.Size
		*totalLines += result.Lines
		fmt.Printf("Processed: %s (Size: %d bytes, Lines: %d)\n", 
			filepath.Base(result.Path), result.Size, result.Lines)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	root := os.Args[1]
	pattern := filepath.Join(root, "*.txt")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	if len(matches) == 0 {
		fmt.Println("No .txt files found in specified directory")
		return
	}

	startTime := time.Now()
	results := make(chan FileResult, len(matches))
	var wg sync.WaitGroup

	for _, file := range matches {
		wg.Add(1)
		go processFile(file, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var totalFiles int
	var totalSize int64
	var totalLines int

	collectResults(results, &totalFiles, &totalSize, &totalLines)

	duration := time.Since(startTime)
	fmt.Printf("\nSummary:\n")
	fmt.Printf("Total files processed: %d\n", totalFiles)
	fmt.Printf("Total size: %d bytes\n", totalSize)
	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Processing time: %v\n", duration)
}