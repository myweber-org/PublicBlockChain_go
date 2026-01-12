
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID    int
    Name  string
    Value float64
}

func processCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
        }

        name := row[1]
        if name == "" {
            return nil, fmt.Errorf("empty name at line %d", lineNumber)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, fmt.Errorf("no valid records found in file")
    }

    return records, nil
}

func validateRecords(records []Record) error {
    seenIDs := make(map[int]bool)
    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("record ID must be positive: %d", record.ID)
        }
        if seenIDs[record.ID] {
            return fmt.Errorf("duplicate ID found: %d", record.ID)
        }
        seenIDs[record.ID] = true

        if record.Value < 0 {
            return fmt.Errorf("negative value not allowed for record %d", record.ID)
        }
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: file_processor <csv_file>")
        os.Exit(1)
    }

    filename := os.Args[1]
    records, err := processCSVFile(filename)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    if err := validateRecords(records); err != nil {
        fmt.Printf("Validation error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %d records\n", len(records))
    totalValue := 0.0
    for _, record := range records {
        totalValue += record.Value
        fmt.Printf("ID: %d, Name: %s, Value: %.2f\n", record.ID, record.Name, record.Value)
    }
    fmt.Printf("Total value: %.2f\n", totalValue)
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
	mu        sync.RWMutex
	stats     ProcessingStats
}

type ProcessingStats struct {
	FilesProcessed int
	TotalBytes     int64
	Errors         int
	StartTime      time.Time
	Duration       time.Duration
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
		batchSize = 100
	}

	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
		stats: ProcessingStats{
			StartTime: time.Now(),
		},
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string, processor func([]byte) ([]byte, error)) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dirPath)
	}

	taskChan := make(chan FileTask, fp.batchSize)
	resultChan := make(chan FileTask, fp.batchSize)
	var wg sync.WaitGroup

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(taskChan, resultChan, &wg, processor)
	}

	go fp.collectResults(resultChan)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Size() > 100*1024*1024 {
			return fmt.Errorf("file too large: %s", path)
		}

		taskChan <- FileTask{Path: path}
		return nil
	})

	close(taskChan)
	wg.Wait()
	close(resultChan)

	fp.mu.Lock()
	fp.stats.Duration = time.Since(fp.stats.StartTime)
	fp.mu.Unlock()

	return err
}

func (fp *FileProcessor) worker(taskChan <-chan FileTask, resultChan chan<- FileTask, wg *sync.WaitGroup, processor func([]byte) ([]byte, error)) {
	defer wg.Done()

	for task := range taskChan {
		content, err := fp.readFile(task.Path)
		if err != nil {
			task.Err = err
			resultChan <- task
			continue
		}

		processed, err := processor(content)
		if err != nil {
			task.Err = err
		} else {
			task.Content = processed
		}

		resultChan <- task
	}
}

func (fp *FileProcessor) readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var content []byte
	buffer := make([]byte, 4096)

	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if n == 0 {
			break
		}

		content = append(content, buffer[:n]...)
	}

	return content, nil
}

func (fp *FileProcessor) collectResults(resultChan <-chan FileTask) {
	for result := range resultChan {
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
}

func (fp *FileProcessor) GetStats() ProcessingStats {
	fp.mu.RLock()
	defer fp.mu.RUnlock()
	return fp.stats
}

func ExampleProcessor(content []byte) ([]byte, error) {
	if len(content) == 0 {
		return nil, errors.New("empty content")
	}

	processed := make([]byte, len(content))
	for i, b := range content {
		processed[i] = b ^ 0xFF
	}

	return processed, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	processor := NewFileProcessor(8, 50)

	fmt.Printf("Processing files in: %s\n", dirPath)

	err := processor.ProcessDirectory(dirPath, ExampleProcessor)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		os.Exit(1)
	}

	stats := processor.GetStats()
	fmt.Printf("\nProcessing completed:\n")
	fmt.Printf("Files processed: %d\n", stats.FilesProcessed)
	fmt.Printf("Total bytes: %d\n", stats.TotalBytes)
	fmt.Printf("Errors: %d\n", stats.Errors)
	fmt.Printf("Duration: %v\n", stats.Duration)
}