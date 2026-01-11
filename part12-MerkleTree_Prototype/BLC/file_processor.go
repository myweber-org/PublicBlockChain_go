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
	mu        sync.RWMutex
	stats     map[string]int
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
		stats:     make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	var filePaths []string
	for _, file := range files {
		if !file.IsDir() {
			filePaths = append(filePaths, filepath.Join(dirPath, file.Name()))
		}
	}

	return fp.processBatch(filePaths)
}

func (fp *FileProcessor) processBatch(paths []string) error {
	var wg sync.WaitGroup
	fileChan := make(chan string, fp.batchSize)
	errChan := make(chan error, fp.workers)

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(&wg, fileChan, errChan)
	}

	for _, path := range paths {
		select {
		case err := <-errChan:
			return err
		default:
			fileChan <- path
		}
	}

	close(fileChan)
	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (fp *FileProcessor) worker(wg *sync.WaitGroup, files <-chan string, errChan chan<- error) {
	defer wg.Done()

	for filePath := range files {
		if err := fp.processFile(filePath); err != nil {
			errChan <- fmt.Errorf("failed to process %s: %w", filePath, err)
			return
		}
	}
}

func (fp *FileProcessor) processFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	wordCount := 0

	for scanner.Scan() {
		lineCount++
		wordCount += countWords(scanner.Text())
	}

	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	fp.mu.Lock()
	fp.stats[path] = wordCount
	fp.mu.Unlock()

	fmt.Printf("Processed %s: %d lines, %d words\n", filepath.Base(path), lineCount, wordCount)
	return nil
}

func countWords(line string) int {
	inWord := false
	count := 0

	for _, ch := range line {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			if !inWord {
				count++
				inWord = true
			}
		} else {
			inWord = false
		}
	}
	return count
}

func (fp *FileProcessor) GetStats() map[string]int {
	fp.mu.RLock()
	defer fp.mu.RUnlock()

	statsCopy := make(map[string]int, len(fp.stats))
	for k, v := range fp.stats {
		statsCopy[k] = v
	}
	return statsCopy
}

func main() {
	processor := NewFileProcessor(4, 10)
	
	start := time.Now()
	if err := processor.ProcessDirectory("."); err != nil {
		fmt.Printf("Processing error: %v\n", err)
		os.Exit(1)
	}
	
	elapsed := time.Since(start)
	stats := processor.GetStats()
	
	fmt.Printf("\nProcessing completed in %v\n", elapsed)
	fmt.Printf("Total files processed: %d\n", len(stats))
}