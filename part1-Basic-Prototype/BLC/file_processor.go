package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type DataRecord struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Processed bool      `json:"processed"`
}

type Processor struct {
	mu          sync.RWMutex
	records     []DataRecord
	workerCount int
	results     chan string
	errors      chan error
}

func NewProcessor(workerCount int) *Processor {
	return &Processor{
		records:     make([]DataRecord, 0),
		workerCount: workerCount,
		results:     make(chan string, 100),
		errors:      make(chan error, 100),
	}
}

func (p *Processor) LoadData(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var records []DataRecord
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&records); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	p.mu.Lock()
	p.records = records
	p.mu.Unlock()

	return nil
}

func (p *Processor) processRecord(record DataRecord) (string, error) {
	time.Sleep(10 * time.Millisecond)

	if record.Value < 0 {
		return "", fmt.Errorf("invalid value for record %d: %f", record.ID, record.Value)
	}

	record.Processed = true
	record.Timestamp = time.Now()

	return fmt.Sprintf("Processed record %d: %s (%.2f)", record.ID, record.Name, record.Value), nil
}

func (p *Processor) worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	p.mu.RLock()
	records := p.records
	p.mu.RUnlock()

	for i := id; i < len(records); i += p.workerCount {
		result, err := p.processRecord(records[i])
		if err != nil {
			p.errors <- err
			continue
		}
		p.results <- result
	}
}

func (p *Processor) Run() {
	var wg sync.WaitGroup

	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		go p.worker(i, &wg)
	}

	wg.Wait()
	close(p.results)
	close(p.errors)
}

func (p *Processor) SaveResults(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	p.mu.RLock()
	defer p.mu.RUnlock()

	if err := encoder.Encode(p.records); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func main() {
	processor := NewProcessor(4)

	if err := processor.LoadData("input.json"); err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	go func() {
		for result := range processor.results {
			fmt.Println(result)
		}
	}()

	go func() {
		for err := range processor.errors {
			log.Printf("Processing error: %v", err)
		}
	}()

	processor.Run()

	if err := processor.SaveResults("output.json"); err != nil {
		log.Fatalf("Failed to save results: %v", err)
	}

	fmt.Println("Processing completed successfully")
}