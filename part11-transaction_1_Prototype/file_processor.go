package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileStats struct {
	Path    string
	Size    int64
	Lines   int
	Words   int
	Chars   int
}

func processFile(path string, wg *sync.WaitGroup, results chan<- FileStats) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", path, err)
		return
	}
	defer file.Close()

	stats := FileStats{Path: path}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info for %s: %v\n", path, err)
		return
	}
	stats.Size = fileInfo.Size()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stats.Lines++
		stats.Chars += len(line)
		
		wordScanner := bufio.NewScanner(bufio.NewReader(filepath.NewReader(line)))
		wordScanner.Split(bufio.ScanWords)
		for wordScanner.Scan() {
			stats.Words++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file %s: %v\n", path, err)
		return
	}

	results <- stats
}

func collectResults(results <-chan FileStats, totalFiles *int, totalStats *FileStats) {
	for stats := range results {
		*totalFiles++
		totalStats.Size += stats.Size
		totalStats.Lines += stats.Lines
		totalStats.Words += stats.Words
		totalStats.Chars += stats.Chars

		fmt.Printf("Processed: %s\n", stats.Path)
		fmt.Printf("  Size: %d bytes, Lines: %d, Words: %d, Characters: %d\n\n",
			stats.Size, stats.Lines, stats.Words, stats.Chars)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_processor.go <file1> [file2] ...")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	results := make(chan FileStats, len(os.Args)-1)
	totalFiles := 0
	totalStats := FileStats{}

	go collectResults(results, &totalFiles, &totalStats)

	for _, filePath := range os.Args[1:] {
		wg.Add(1)
		go processFile(filePath, &wg, results)
	}

	wg.Wait()
	close(results)

	fmt.Printf("\n=== SUMMARY ===\n")
	fmt.Printf("Total files processed: %d\n", totalFiles)
	fmt.Printf("Combined size: %d bytes\n", totalStats.Size)
	fmt.Printf("Total lines: %d\n", totalStats.Lines)
	fmt.Printf("Total words: %d\n", totalStats.Words)
	fmt.Printf("Total characters: %d\n", totalStats.Chars)
}package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type DataRecord struct {
	ID        string    `json:"id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Processed bool      `json:"processed"`
}

type Processor struct {
	records []DataRecord
	mu      sync.Mutex
	wg      sync.WaitGroup
}

func NewProcessor() *Processor {
	return &Processor{
		records: make([]DataRecord, 0),
	}
}

func (p *Processor) AddRecord(record DataRecord) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.records = append(p.records, record)
}

func (p *Processor) ProcessRecord(index int) {
	defer p.wg.Done()

	if index < 0 || index >= len(p.records) {
		log.Printf("Invalid index %d for processing", index)
		return
	}

	p.mu.Lock()
	record := &p.records[index]
	p.mu.Unlock()

	time.Sleep(50 * time.Millisecond)

	p.mu.Lock()
	record.Processed = true
	record.Value = record.Value * 1.1
	p.mu.Unlock()

	log.Printf("Processed record %s with new value %.2f", record.ID, record.Value)
}

func (p *Processor) ProcessAll() {
	p.wg.Add(len(p.records))
	for i := range p.records {
		go p.ProcessRecord(i)
	}
	p.wg.Wait()
}

func (p *Processor) GenerateReport() {
	p.mu.Lock()
	defer p.mu.Unlock()

	processedCount := 0
	totalValue := 0.0

	for _, record := range p.records {
		if record.Processed {
			processedCount++
			totalValue += record.Value
		}
	}

	fmt.Printf("Processing Report:\n")
	fmt.Printf("Total records: %d\n", len(p.records))
	fmt.Printf("Processed records: %d\n", processedCount)
	fmt.Printf("Average value: %.2f\n", totalValue/float64(processedCount))
}

func (p *Processor) SaveToFile(filename string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(p.records); err != nil {
		return fmt.Errorf("failed to encode records: %w", err)
	}

	return nil
}

func main() {
	processor := NewProcessor()

	for i := 0; i < 10; i++ {
		record := DataRecord{
			ID:        fmt.Sprintf("REC-%03d", i+1),
			Value:     float64(i) * 10.5,
			Timestamp: time.Now(),
			Processed: false,
		}
		processor.AddRecord(record)
	}

	fmt.Println("Starting concurrent processing...")
	startTime := time.Now()
	processor.ProcessAll()
	duration := time.Since(startTime)

	processor.GenerateReport()
	fmt.Printf("Processing completed in %v\n", duration)

	if err := processor.SaveToFile("processed_data.json"); err != nil {
		log.Fatalf("Failed to save data: %v", err)
	}

	fmt.Println("Data saved to processed_data.json")
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
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Valid     bool      `json:"valid"`
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

func (p *Processor) AddRecord(id string, value float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	record := DataRecord{
		ID:        id,
		Timestamp: time.Now().UTC(),
		Value:     value,
		Valid:     value >= 0 && value <= 100,
	}

	p.records = append(p.records, record)
}

func (p *Processor) ValidateRecords() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var wg sync.WaitGroup
	results := make(chan string, len(p.records))

	for _, record := range p.records {
		wg.Add(1)
		go func(r DataRecord) {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			if r.Valid {
				results <- fmt.Sprintf("Record %s: PASS", r.ID)
			} else {
				results <- fmt.Sprintf("Record %s: FAIL", r.ID)
			}
		}(record)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}

func (p *Processor) GenerateReport() ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	summary := struct {
		Total     int `json:"total_records"`
		Valid     int `json:"valid_records"`
		Invalid   int `json:"invalid_records"`
		Timestamp string `json:"generated_at"`
	}{
		Total:     len(p.records),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	for _, record := range p.records {
		if record.Valid {
			summary.Valid++
		} else {
			summary.Invalid++
		}
	}

	return json.MarshalIndent(summary, "", "  ")
}

func main() {
	processor := NewProcessor()

	processor.AddRecord("A001", 45.7)
	processor.AddRecord("A002", 102.3)
	processor.AddRecord("A003", 78.9)
	processor.AddRecord("A004", -5.2)

	fmt.Println("Validating records...")
	processor.ValidateRecords()

	report, err := processor.GenerateReport()
	if err != nil {
		fmt.Printf("Error generating report: %v\n", err)
		return
	}

	fmt.Println("\nProcessing Report:")
	fmt.Println(string(report))
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
	mu          sync.RWMutex
	processed   map[string]bool
	maxWorkers  int
	fileChannel chan string
	wg          sync.WaitGroup
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		processed:   make(map[string]bool),
		maxWorkers:  workers,
		fileChannel: make(chan string, 100),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	if !fp.markProcessing(path) {
		return errors.New("file already being processed")
	}
	defer fp.markComplete(path)

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	wordCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		wordCount += countWords(line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	fp.mu.Lock()
	fmt.Printf("Processed %s: lines=%d words=%d\n", filepath.Base(path), lineCount, wordCount)
	fp.mu.Unlock()
	return nil
}

func (fp *FileProcessor) markProcessing(path string) bool {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	if fp.processed[path] {
		return false
	}
	fp.processed[path] = true
	return true
}

func (fp *FileProcessor) markComplete(path string) {
	fp.mu.Lock()
	delete(fp.processed, path)
	fp.mu.Unlock()
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

func (fp *FileProcessor) worker(id int) {
	defer fp.wg.Done()
	for path := range fp.fileChannel {
		start := time.Now()
		err := fp.ProcessFile(path)
		duration := time.Since(start)
		if err != nil {
			fmt.Printf("Worker %d error processing %s: %v\n", id, path, err)
		} else {
			fmt.Printf("Worker %d completed %s in %v\n", id, filepath.Base(path), duration)
		}
	}
}

func (fp *FileProcessor) StartWorkers() {
	for i := 0; i < fp.maxWorkers; i++ {
		fp.wg.Add(1)
		go fp.worker(i + 1)
	}
}

func (fp *FileProcessor) AddFile(path string) {
	fp.fileChannel <- path
}

func (fp *FileProcessor) WaitAndClose() {
	close(fp.fileChannel)
	fp.wg.Wait()
}

func main() {
	processor := NewFileProcessor(4)
	processor.StartWorkers()

	files := []string{"test1.txt", "test2.txt", "test3.txt"}
	for _, file := range files {
		processor.AddFile(file)
	}

	processor.WaitAndClose()
	fmt.Println("All files processed")
}