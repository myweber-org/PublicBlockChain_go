
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
	TotalBytes     int64
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
		stats: ProcessingStats{
			StartTime: time.Now(),
		},
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	filePaths := make([]string, 0)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Mode().IsRegular() {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	return fp.processFiles(filePaths)
}

func (fp *FileProcessor) processFiles(filePaths []string) error {
	taskChan := make(chan FileTask, fp.batchSize)
	results := make(chan FileTask, len(filePaths))
	var wg sync.WaitGroup

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(taskChan, results, &wg)
	}

	go func() {
		for _, path := range filePaths {
			taskChan <- FileTask{Path: path}
		}
		close(taskChan)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fp.mu.Lock()
		fp.stats.FilesProcessed++
		if result.Err != nil {
			fp.stats.Errors++
			fmt.Printf("Error processing %s: %v\n", result.Path, result.Err)
		} else {
			fp.stats.TotalBytes += int64(len(result.Content))
		}
		fp.mu.Unlock()
	}

	fp.stats.EndTime = time.Now()
	return nil
}

func (fp *FileProcessor) worker(tasks <-chan FileTask, results chan<- FileTask, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		content, err := fp.readFile(task.Path)
		task.Content = content
		task.Err = err
		results <- task
	}
}

func (fp *FileProcessor) readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var content []byte
	buffer := make([]byte, 4096)

	for {
		n, err := reader.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		if n == 0 {
			break
		}
		content = append(content, buffer[:n]...)
	}

	return content, nil
}

func (fp *FileProcessor) GetStats() ProcessingStats {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.stats
}

func (fp *FileProcessor) PrintStats() {
	stats := fp.GetStats()
	duration := stats.EndTime.Sub(stats.StartTime)

	fmt.Println("\n=== Processing Statistics ===")
	fmt.Printf("Files processed: %d\n", stats.FilesProcessed)
	fmt.Printf("Total bytes: %d\n", stats.TotalBytes)
	fmt.Printf("Errors: %d\n", stats.Errors)
	fmt.Printf("Duration: %v\n", duration)
	if duration > 0 {
		fmt.Printf("Throughput: %.2f MB/s\n", 
			float64(stats.TotalBytes)/(1024*1024)/duration.Seconds())
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory_path>")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	processor := NewFileProcessor(4, 20)

	fmt.Printf("Processing directory: %s\n", dirPath)
	
	if err := processor.ProcessDirectory(dirPath); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}

	processor.PrintStats()
}