package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu       sync.Mutex
	fileList []string
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		fileList: make([]string, 0),
	}
}

func (fp *FileProcessor) ScanDirectory(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fp.mu.Lock()
			fp.fileList = append(fp.fileList, path)
			fp.mu.Unlock()
		}
		return nil
	})
}

func (fp *FileProcessor) ProcessFiles(workerCount int) {
	var wg sync.WaitGroup
	fileChan := make(chan string)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for filePath := range fileChan {
				fp.processSingleFile(filePath, workerID)
			}
		}(i)
	}

	for _, file := range fp.fileList {
		fileChan <- file
	}
	close(fileChan)
	wg.Wait()
}

func (fp *FileProcessor) processSingleFile(filePath string, workerID int) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Worker %d: Failed to open %s: %v\n", workerID, filePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Worker %d: Error reading %s: %v\n", workerID, filePath, err)
		return
	}

	fmt.Printf("Worker %d: Processed %s - %d lines\n", workerID, filePath, lineCount)
}

func (fp *FileProcessor) GetFileCount() int {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return len(fp.fileList)
}

func main() {
	processor := NewFileProcessor()
	
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	err := processor.ScanDirectory(dirPath)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d files\n", processor.GetFileCount())
	processor.ProcessFiles(4)
}
package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type DataRecord struct {
	ID        string    `json:"id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Valid     bool      `json:"valid"`
}

type Processor struct {
	mu      sync.RWMutex
	records map[string]DataRecord
}

func NewProcessor() *Processor {
	return &Processor{
		records: make(map[string]DataRecord),
	}
}

func (p *Processor) AddRecord(record DataRecord) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	record.Valid = record.Value >= 0 && record.Value <= 1000
	p.records[record.ID] = record
}

func (p *Processor) ValidateRecords() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	var invalidIDs []string
	for id, record := range p.records {
		if !record.Valid {
			invalidIDs = append(invalidIDs, id)
		}
	}
	return invalidIDs
}

func (p *Processor) TransformRecords(multiplier float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	for id, record := range p.records {
		if record.Valid {
			record.Value *= multiplier
			p.records[id] = record
		}
	}
}

func (p *Processor) ExportJSON() ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	return json.MarshalIndent(p.records, "", "  ")
}

func worker(id int, jobs <-chan DataRecord, results chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for record := range jobs {
		time.Sleep(10 * time.Millisecond)
		results <- record.Valid
	}
}

func main() {
	processor := NewProcessor()
	
	records := []DataRecord{
		{ID: "A1", Value: 450.5, Timestamp: time.Now()},
		{ID: "B2", Value: -10.0, Timestamp: time.Now()},
		{ID: "C3", Value: 1200.0, Timestamp: time.Now()},
		{ID: "D4", Value: 750.3, Timestamp: time.Now()},
	}
	
	for _, record := range records {
		processor.AddRecord(record)
	}
	
	fmt.Println("Invalid records:", processor.ValidateRecords())
	
	processor.TransformRecords(1.1)
	
	jsonData, err := processor.ExportJSON()
	if err != nil {
		fmt.Println("Export error:", err)
		return
	}
	
	fmt.Println("Processed data:")
	fmt.Println(string(jsonData))
	
	numWorkers := 3
	jobs := make(chan DataRecord, len(records))
	results := make(chan bool, len(records))
	var wg sync.WaitGroup
	
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}
	
	for _, record := range records {
		jobs <- record
	}
	close(jobs)
	
	wg.Wait()
	close(results)
	
	validCount := 0
	for result := range results {
		if result {
			validCount++
		}
	}
	
	fmt.Printf("Concurrent validation: %d/%d records valid\n", validCount, len(records))
}package main

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
	BatchSize   int
	WorkerCount int
	mu          sync.Mutex
	processed   int
}

func NewFileProcessor(batchSize, workerCount int) *FileProcessor {
	return &FileProcessor{
		BatchSize:   batchSize,
		WorkerCount: workerCount,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	fileChan := make(chan string, len(paths))
	resultChan := make(chan bool, len(paths))
	var wg sync.WaitGroup

	for i := 0; i < fp.WorkerCount; i++ {
		wg.Add(1)
		go fp.worker(fileChan, resultChan, &wg)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	wg.Wait()
	close(resultChan)

	successCount := 0
	for result := range resultChan {
		if result {
			successCount++
		}
	}

	fp.mu.Lock()
	fp.processed += successCount
	fp.mu.Unlock()

	fmt.Printf("Processed %d files successfully\n", successCount)
	return nil
}

func (fp *FileProcessor) worker(files <-chan string, results chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range files {
		success := fp.processSingleFile(file)
		results <- success
	}
}

func (fp *FileProcessor) processSingleFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", path, err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		if lineCount >= fp.BatchSize {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file %s: %v\n", path, err)
		return false
	}

	time.Sleep(10 * time.Millisecond)
	fmt.Printf("Processed %s (%d lines)\n", filepath.Base(path), lineCount)
	return true
}

func (fp *FileProcessor) GetProcessedCount() int {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.processed
}

func main() {
	processor := NewFileProcessor(100, 4)

	files := []string{
		"data/file1.txt",
		"data/file2.txt",
		"data/file3.txt",
		"data/file4.txt",
	}

	if err := processor.ProcessFiles(files); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		return
	}

	fmt.Printf("Total files processed: %d\n", processor.GetProcessedCount())
}