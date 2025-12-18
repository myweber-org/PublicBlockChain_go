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