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
	Workers    int
	BatchSize  int
	ResultChan chan ProcessResult
	ErrorChan  chan error
}

type ProcessResult struct {
	Filename string
	Lines    int
	Duration time.Duration
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		Workers:    workers,
		BatchSize:  batchSize,
		ResultChan: make(chan ProcessResult, 100),
		ErrorChan:  make(chan error, 100),
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string) {
	var wg sync.WaitGroup
	fileChan := make(chan string, len(paths))

	for i := 0; i < fp.Workers; i++ {
		wg.Add(1)
		go fp.worker(&wg, fileChan)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	wg.Wait()
	close(fp.ResultChan)
	close(fp.ErrorChan)
}

func (fp *FileProcessor) worker(wg *sync.WaitGroup, files <-chan string) {
	defer wg.Done()

	for file := range files {
		result, err := fp.processSingleFile(file)
		if err != nil {
			fp.ErrorChan <- fmt.Errorf("file %s: %w", file, err)
			continue
		}
		fp.ResultChan <- result
	}
}

func (fp *FileProcessor) processSingleFile(path string) (ProcessResult, error) {
	start := time.Now()

	file, err := os.Open(path)
	if err != nil {
		return ProcessResult{}, err
	}
	defer file.Close()

	lineCount := 0
	scanner := bufio.NewScanner(file)
	batch := make([]string, 0, fp.BatchSize)

	for scanner.Scan() {
		batch = append(batch, scanner.Text())
		if len(batch) >= fp.BatchSize {
			if err := fp.processBatch(batch); err != nil {
				return ProcessResult{}, err
			}
			lineCount += len(batch)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := fp.processBatch(batch); err != nil {
			return ProcessResult{}, err
		}
		lineCount += len(batch)
	}

	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		return ProcessResult{}, err
	}

	duration := time.Since(start)
	return ProcessResult{
		Filename: filepath.Base(path),
		Lines:    lineCount,
		Duration: duration,
	}, nil
}

func (fp *FileProcessor) processBatch(lines []string) error {
	if len(lines) == 0 {
		return nil
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}

func main() {
	processor := NewFileProcessor(4, 100)

	files := []string{
		"data/file1.txt",
		"data/file2.txt",
		"data/file3.txt",
	}

	go processor.ProcessFiles(files)

	for result := range processor.ResultChan {
		fmt.Printf("Processed %s: %d lines in %v\n",
			result.Filename, result.Lines, result.Duration)
	}

	for err := range processor.ErrorChan {
		fmt.Printf("Error: %v\n", err)
	}
}package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

type FileProcessor struct {
	workers int
	results chan string
	wg      sync.WaitGroup
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		workers: workers,
		results: make(chan string, 100),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make(chan string, fp.workers)

	for i := 0; i < fp.workers; i++ {
		fp.wg.Add(1)
		go fp.worker(lines)
	}

	go func() {
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	go func() {
		fp.wg.Wait()
		close(fp.results)
	}()

	return scanner.Err()
}

func (fp *FileProcessor) worker(lines <-chan string) {
	defer fp.wg.Done()
	for line := range lines {
		processed := fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), line)
		fp.results <- processed
	}
}

func (fp *FileProcessor) GetResults() <-chan string {
	return fp.results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <filename>")
		os.Exit(1)
	}

	processor := NewFileProcessor(4)
	
	if err := processor.ProcessFile(os.Args[1]); err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	for result := range processor.GetResults() {
		fmt.Println(result)
	}
}