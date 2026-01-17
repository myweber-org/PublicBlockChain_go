
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type DataRecord struct {
	ID        string    `json:"id"`
	Value     int       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
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

func (p *Processor) AddRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value must be non-negative")
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.records = append(p.records, record)
	return nil
}

func (p *Processor) ValidateRecords() {
	p.mu.RLock()
	records := make([]DataRecord, len(p.records))
	copy(records, p.records)
	p.mu.RUnlock()

	var wg sync.WaitGroup
	results := make(chan DataRecord, len(records))

	for _, record := range records {
		wg.Add(1)
		go func(r DataRecord) {
			defer wg.Done()
			r.Valid = r.Value > 0 && r.Timestamp.Before(time.Now())
			results <- r
		}(record)
	}

	wg.Wait()
	close(results)

	p.mu.Lock()
	p.records = make([]DataRecord, 0)
	for result := range results {
		p.records = append(p.records, result)
	}
	p.mu.Unlock()
}

func (p *Processor) ExportJSON() (string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.records) == 0 {
		return "", errors.New("no records to export")
	}

	data, err := json.MarshalIndent(p.records, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal records: %w", err)
	}
	return string(data), nil
}

func (p *Processor) GetStats() (int, int) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	total := len(p.records)
	validCount := 0
	for _, record := range p.records {
		if record.Valid {
			validCount++
		}
	}
	return total, validCount
}

func main() {
	processor := NewProcessor()

	records := []DataRecord{
		{ID: "rec1", Value: 42, Timestamp: time.Now().Add(-1 * time.Hour)},
		{ID: "rec2", Value: 0, Timestamp: time.Now().Add(2 * time.Hour)},
		{ID: "rec3", Value: -5, Timestamp: time.Now().Add(-30 * time.Minute)},
	}

	for _, record := range records {
		if err := processor.AddRecord(record); err != nil {
			fmt.Printf("Failed to add record %s: %v\n", record.ID, err)
		}
	}

	processor.ValidateRecords()

	total, valid := processor.GetStats()
	fmt.Printf("Processed %d records, %d valid\n", total, valid)

	jsonOutput, err := processor.ExportJSON()
	if err != nil {
		fmt.Printf("Export failed: %v\n", err)
	} else {
		fmt.Println("Exported data:")
		fmt.Println(jsonOutput)
	}
}
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

func (p *Processor) AddRecord(id int, value string) error {
	if id <= 0 {
		return errors.New("invalid record ID")
	}
	if value == "" {
		return errors.New("record value cannot be empty")
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
	var wg sync.WaitGroup
	p.mu.RLock()
	records := make([]DataRecord, len(p.records))
	copy(records, p.records)
	p.mu.RUnlock()

	for i := range records {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			p.validateRecord(&records[idx])
		}(i)
	}
	wg.Wait()

	p.mu.Lock()
	for i := range p.records {
		p.records[i].Valid = records[i].Valid
	}
	p.mu.Unlock()
}

func (p *Processor) validateRecord(record *DataRecord) {
	if record.ID <= 0 || record.Value == "" {
		record.Valid = false
		return
	}

	if time.Since(record.Timestamp) > 24*time.Hour {
		record.Valid = false
		return
	}

	record.Valid = true
}

func (p *Processor) GetValidRecords() []DataRecord {
	p.mu.RLock()
	defer p.mu.RUnlock()

	validRecords := make([]DataRecord, 0)
	for _, record := range p.records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func (p *Processor) ProcessBatch(ids []int, values []string) error {
	if len(ids) != len(values) {
		return errors.New("ids and values length mismatch")
	}

	var wg sync.WaitGroup
	errorsChan := make(chan error, len(ids))

	for i := 0; i < len(ids); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if err := p.AddRecord(ids[idx], values[idx]); err != nil {
				errorsChan <- fmt.Errorf("record %d: %v", ids[idx], err)
			}
		}(i)
	}

	wg.Wait()
	close(errorsChan)

	if len(errorsChan) > 0 {
		var errMsg string
		for err := range errorsChan {
			if errMsg != "" {
				errMsg += "; "
			}
			errMsg += err.Error()
		}
		return errors.New(errMsg)
	}

	return nil
}

func main() {
	processor := NewProcessor()

	ids := []int{1, 2, 3, 4, 5}
	values := []string{"alpha", "beta", "gamma", "delta", "epsilon"}

	if err := processor.ProcessBatch(ids, values); err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	processor.ValidateRecords()
	validRecords := processor.GetValidRecords()

	fmt.Printf("Successfully processed %d valid records\n", len(validRecords))
	for _, record := range validRecords {
		fmt.Printf("ID: %d, Value: %s, Time: %s\n",
			record.ID, record.Value, record.Timestamp.Format(time.RFC3339))
	}
}