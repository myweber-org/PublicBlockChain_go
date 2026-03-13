
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	Workers   int
	BatchSize int
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		Workers:   workers,
		BatchSize: batchSize,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string, processor func(string) error) []error {
	var wg sync.WaitGroup
	errorChan := make(chan error, len(paths))
	pathChan := make(chan string, len(paths))

	for i := 0; i < fp.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathChan {
				if err := processor(path); err != nil {
					errorChan <- fmt.Errorf("processing %s: %w", path, err)
				}
			}
		}()
	}

	for _, path := range paths {
		pathChan <- path
	}
	close(pathChan)
	wg.Wait()
	close(errorChan)

	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}
	return errors
}

func (fp *FileProcessor) ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (fp *FileProcessor) BatchProcess(dir string, pattern string) error {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return err
	}

	batches := fp.createBatches(matches)
	for _, batch := range batches {
		errors := fp.ProcessFiles(batch, func(path string) error {
			content, err := fp.ReadLines(path)
			if err != nil {
				return err
			}
			fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), len(content))
			return nil
		})

		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Printf("Error: %v\n", err)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func (fp *FileProcessor) createBatches(items []string) [][]string {
	var batches [][]string
	for i := 0; i < len(items); i += fp.BatchSize {
		end := i + fp.BatchSize
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}
	return batches
}

func main() {
	processor := NewFileProcessor(4, 10)
	if err := processor.BatchProcess(".", "*.txt"); err != nil {
		fmt.Printf("Batch processing failed: %v\n", err)
	}
}package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Record struct {
	ID      int
	Name    string
	Value   float64
	Valid   bool
}

func processFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]Record, 0)
	var wg sync.WaitGroup
	recordChan := make(chan Record, 100)
	errorChan := make(chan error, 10)
	done := make(chan bool)

	go func() {
		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errorChan <- fmt.Errorf("csv read error: %w", err)
				continue
			}

			wg.Add(1)
			go func(data []string) {
				defer wg.Done()
				record, err := validateAndTransform(data)
				if err != nil {
					errorChan <- err
					return
				}
				recordChan <- record
			}(row)
		}
		wg.Wait()
		close(recordChan)
		close(errorChan)
		done <- true
	}()

	for {
		select {
		case record := <-recordChan:
			records = append(records, record)
		case err := <-errorChan:
			fmt.Printf("Processing error: %v\n", err)
		case <-done:
			return records, nil
		}
	}
}

func validateAndTransform(data []string) (Record, error) {
	if len(data) != 4 {
		return Record{}, errors.New("invalid data length")
	}

	id, err := strconv.Atoi(strings.TrimSpace(data[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID format: %w", err)
	}

	name := strings.TrimSpace(data[1])
	if name == "" {
		return Record{}, errors.New("name cannot be empty")
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(data[2]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid value format: %w", err)
	}

	valid := strings.ToLower(strings.TrimSpace(data[3])) == "true"

	return Record{
		ID:    id,
		Name:  name,
		Value: value,
		Valid: valid,
	}, nil
}

func aggregateResults(records []Record) map[string]float64 {
	aggregation := make(map[string]float64)
	for _, record := range records {
		if record.Valid {
			aggregation[record.Name] += record.Value
		}
	}
	return aggregation
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <filename>")
		os.Exit(1)
	}

	records, err := processFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	results := aggregateResults(records)
	for name, total := range results {
		fmt.Printf("%s: %.2f\n", name, total)
	}

	fmt.Printf("Processed %d records successfully\n", len(records))
}