
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileStats struct {
	Path     string
	Size     int64
	Lines    int
	Modified time.Time
}

func processFile(path string, wg *sync.WaitGroup, results chan<- FileStats) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", path, err)
		return
	}
	defer file.Close()

	stats := FileStats{Path: path}

	if fileInfo, err := file.Stat(); err == nil {
		stats.Size = fileInfo.Size()
		stats.Modified = fileInfo.ModTime()
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	stats.Lines = lineCount

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning %s: %v\n", path, err)
		return
	}

	results <- stats
}

func collectFiles(dir string, patterns []string) ([]string, error) {
	var files []string

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return nil, err
		}
		files = append(files, matches...)
	}

	return files, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory> [patterns...]")
		fmt.Println("Default patterns: *.txt, *.go, *.md")
		os.Exit(1)
	}

	dir := os.Args[1]
	patterns := []string{"*.txt", "*.go", "*.md"}
	if len(os.Args) > 2 {
		patterns = os.Args[2:]
	}

	files, err := collectFiles(dir, patterns)
	if err != nil {
		fmt.Printf("Error collecting files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No files found matching patterns")
		return
	}

	var wg sync.WaitGroup
	results := make(chan FileStats, len(files))

	fmt.Printf("Processing %d files...\n", len(files))

	for _, file := range files {
		wg.Add(1)
		go processFile(file, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalSize := int64(0)
	totalLines := 0
	fileCount := 0

	for stats := range results {
		fmt.Printf("%s: %d bytes, %d lines, modified %s\n",
			filepath.Base(stats.Path),
			stats.Size,
			stats.Lines,
			stats.Modified.Format("2006-01-02 15:04:05"))

		totalSize += stats.Size
		totalLines += stats.Lines
		fileCount++
	}

	fmt.Printf("\nSummary: Processed %d files, total size: %d bytes, total lines: %d\n",
		fileCount, totalSize, totalLines)
}
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type DataRecord struct {
	ID        string
	Value     int
	Timestamp time.Time
	Valid     bool
}

type Processor struct {
	records []DataRecord
	mu      sync.RWMutex
}

func NewProcessor() *Processor {
	return &Processor{
		records: make([]DataRecord, 0),
	}
}

func (p *Processor) AddRecord(id string, value int) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	if value < 0 {
		return errors.New("value must be non-negative")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	record := DataRecord{
		ID:        id,
		Value:     value,
		Timestamp: time.Now(),
		Valid:     true,
	}
	p.records = append(p.records, record)
	return nil
}

func (p *Processor) ValidateRecords() {
	p.mu.RLock()
	records := make([]DataRecord, len(p.records))
	copy(records, p.records)
	p.mu.RUnlock()

	var wg sync.WaitGroup
	for i := range records {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			p.validateRecord(&records[idx])
		}(i)
	}
	wg.Wait()

	p.mu.Lock()
	p.records = records
	p.mu.Unlock()
}

func (p *Processor) validateRecord(record *DataRecord) {
	if record.Value > 1000 {
		record.Valid = false
	}
	if time.Since(record.Timestamp) > 24*time.Hour {
		record.Valid = false
	}
}

func (p *Processor) GetValidCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	count := 0
	for _, record := range p.records {
		if record.Valid {
			count++
		}
	}
	return count
}

func main() {
	processor := NewProcessor()

	records := []struct {
		id    string
		value int
	}{
		{"A001", 150},
		{"A002", 2500},
		{"A003", 75},
		{"A004", -5},
		{"", 100},
	}

	for _, r := range records {
		err := processor.AddRecord(r.id, r.value)
		if err != nil {
			fmt.Printf("Failed to add record %s: %v\n", r.id, err)
		}
	}

	processor.ValidateRecords()
	validCount := processor.GetValidCount()
	fmt.Printf("Valid records: %d out of %d\n", validCount, len(records))
}package main

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

	var wg sync.WaitGroup
	fileChan := make(chan string, fp.batchSize)

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(&wg, fileChan)
	}

	for _, file := range files {
		if !file.IsDir() {
			fullPath := filepath.Join(dirPath, file.Name())
			fileChan <- fullPath
		}
	}

	close(fileChan)
	wg.Wait()

	return fp.printStats()
}

func (fp *FileProcessor) worker(wg *sync.WaitGroup, fileChan <-chan string) {
	defer wg.Done()

	for filePath := range fileChan {
		if err := fp.processFile(filePath); err != nil {
			fmt.Printf("Error processing %s: %v\n", filePath, err)
		}
	}
}

func (fp *FileProcessor) processFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	wordCount := 0

	for scanner.Scan() {
		lineCount++
		words := splitWords(scanner.Text())
		wordCount += len(words)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	fp.mu.Lock()
	fp.stats["files"]++
	fp.stats["lines"] += lineCount
	fp.stats["words"] += wordCount
	fp.mu.Unlock()

	fmt.Printf("Processed: %s (lines: %d, words: %d)\n", filepath.Base(filePath), lineCount, wordCount)
	return nil
}

func splitWords(text string) []string {
	var words []string
	wordStart := -1

	for i, r := range text {
		if isWordRune(r) {
			if wordStart == -1 {
				wordStart = i
			}
		} else {
			if wordStart != -1 {
				words = append(words, text[wordStart:i])
				wordStart = -1
			}
		}
	}

	if wordStart != -1 {
		words = append(words, text[wordStart:])
	}

	return words
}

func isWordRune(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func (fp *FileProcessor) printStats() error {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	if len(fp.stats) == 0 {
		return errors.New("no files processed")
	}

	fmt.Println("\nProcessing Summary:")
	fmt.Printf("Total files: %d\n", fp.stats["files"])
	fmt.Printf("Total lines: %d\n", fp.stats["lines"])
	fmt.Printf("Total words: %d\n", fp.stats["words"])
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	processor := NewFileProcessor(4, 10)

	startTime := time.Now()
	if err := processor.ProcessDirectory(dirPath); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}

	duration := time.Since(startTime)
	fmt.Printf("\nProcessing completed in %v\n", duration)
}