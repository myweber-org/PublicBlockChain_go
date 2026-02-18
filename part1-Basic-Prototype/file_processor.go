package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Server string `json:"server"`
	Port   int    `json:"port"`
	Debug  bool   `json:"debug"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return &config, nil
}

func saveConfig(filename string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func main() {
	config := &Config{
		Server: "api.example.com",
		Port:   8080,
		Debug:  true,
	}

	if err := saveConfig("config.json", config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Config saved successfully")

	loadedConfig, err := loadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded config: %+v\n", loadedConfig)
}package main

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
	fileChan := make(chan string, len(fp.fileList))

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
	return len(fp.fileList)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory_path>")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	processor := NewFileProcessor()

	fmt.Printf("Scanning directory: %s\n", dirPath)
	err := processor.ScanDirectory(dirPath)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d files\n", processor.GetFileCount())
	processor.ProcessFiles(4)
	fmt.Println("Processing completed")
}