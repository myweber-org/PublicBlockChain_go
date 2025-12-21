package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type DataRecord struct {
	ID        int
	Value     string
	Validated bool
	Timestamp time.Time
}

type Processor struct {
	records []DataRecord
	mu      sync.RWMutex
	wg      sync.WaitGroup
}

func NewProcessor() *Processor {
	return &Processor{
		records: make([]DataRecord, 0),
	}
}

func (p *Processor) AddRecord(id int, value string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	record := DataRecord{
		ID:        id,
		Value:     value,
		Validated: false,
		Timestamp: time.Now(),
	}
	p.records = append(p.records, record)
}

func (p *Processor) ValidateRecord(index int) error {
	if index < 0 || index >= len(p.records) {
		return errors.New("index out of bounds")
	}
	
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.records[index].Value == "" {
		return errors.New("empty value detected")
	}
	
	p.records[index].Validated = true
	return nil
}

func (p *Processor) ProcessBatch(start, end int) {
	defer p.wg.Done()
	
	for i := start; i < end && i < len(p.records); i++ {
		err := p.ValidateRecord(i)
		if err != nil {
			fmt.Printf("Validation failed for record %d: %v\n", i, err)
			continue
		}
		fmt.Printf("Record %d validated successfully\n", i)
	}
}

func (p *Processor) RunConcurrentValidation(workers int) {
	batchSize := len(p.records) / workers
	
	for i := 0; i < workers; i++ {
		start := i * batchSize
		end := start + batchSize
		
		if i == workers-1 {
			end = len(p.records)
		}
		
		p.wg.Add(1)
		go p.ProcessBatch(start, end)
	}
	
	p.wg.Wait()
}

func (p *Processor) GetValidatedCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	count := 0
	for _, record := range p.records {
		if record.Validated {
			count++
		}
	}
	return count
}

func main() {
	processor := NewProcessor()
	
	for i := 1; i <= 20; i++ {
		value := fmt.Sprintf("data-%d", i)
		if i%7 == 0 {
			value = ""
		}
		processor.AddRecord(i, value)
	}
	
	fmt.Printf("Total records: %d\n", len(processor.records))
	
	processor.RunConcurrentValidation(4)
	
	validated := processor.GetValidatedCount()
	fmt.Printf("Successfully validated records: %d/%d\n", validated, len(processor.records))
}