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