package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Server string `json:"server"`
	Port   int    `json:"port"`
	Debug  bool   `json:"debug"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(filename string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func main() {
	config := &Config{
		Server: "api.example.com",
		Port:   8080,
		Debug:  true,
	}

	err := SaveConfig("config.json", config)
	if err != nil {
		log.Fatal(err)
	}

	loadedConfig, err := LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Loaded config: %+v", loadedConfig)

	os.Remove("config.json")
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
	stats     map[string]int
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
		stats:     make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, fp.batchSize)

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(i, fileChan, &wg)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	wg.Wait()
	return nil
}

func (fp *FileProcessor) worker(id int, files <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range files {
		start := time.Now()
		err := fp.processSingleFile(file)
		duration := time.Since(start)

		fp.mu.Lock()
		if err != nil {
			fp.stats["failed"]++
			fmt.Printf("Worker %d: Failed to process %s: %v\n", id, file, err)
		} else {
			fp.stats["processed"]++
			fmt.Printf("Worker %d: Processed %s in %v\n", id, file, duration)
		}
		fp.mu.Unlock()
	}
}

func (fp *FileProcessor) processSingleFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := filepath.Ext(path)
	switch ext {
	case ".txt":
		return fp.processTextFile(file)
	case ".log":
		return fp.processLogFile(file)
	default:
		return fp.processGenericFile(file)
	}
}

func (fp *FileProcessor) processTextFile(file io.Reader) error {
	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		_ = scanner.Text()
	}

	fp.mu.Lock()
	fp.stats["lines_read"] += lineCount
	fp.mu.Unlock()

	return scanner.Err()
}

func (fp *FileProcessor) processLogFile(file io.Reader) error {
	return fp.processTextFile(file)
}

func (fp *FileProcessor) processGenericFile(file io.Reader) error {
	buffer := make([]byte, 1024)
	totalBytes := 0

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		totalBytes += n
	}

	fp.mu.Lock()
	fp.stats["bytes_read"] += totalBytes
	fp.mu.Unlock()

	return nil
}

func (fp *FileProcessor) GetStats() map[string]int {
	fp.mu.RLock()
	defer fp.mu.RUnlock()

	statsCopy := make(map[string]int)
	for k, v := range fp.stats {
		statsCopy[k] = v
	}
	return statsCopy
}

func main() {
	processor := NewFileProcessor(4, 10)

	files := []string{
		"file1.txt",
		"file2.log",
		"file3.txt",
		"file4.dat",
	}

	err := processor.ProcessFiles(files)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
	}

	stats := processor.GetStats()
	fmt.Println("\nProcessing statistics:")
	for key, value := range stats {
		fmt.Printf("%s: %d\n", key, value)
	}
}