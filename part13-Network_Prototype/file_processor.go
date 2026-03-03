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
	results  map[string]int
	wg       sync.WaitGroup
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		results: make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	fp.wg.Add(1)
	go func() {
		defer fp.wg.Done()
		
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", path, err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineCount := 0
		for scanner.Scan() {
			lineCount++
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error scanning file %s: %v\n", path, err)
			return
		}

		fp.mu.Lock()
		fp.results[path] = lineCount
		fp.mu.Unlock()
	}()

	return nil
}

func (fp *FileProcessor) Wait() {
	fp.wg.Wait()
}

func (fp *FileProcessor) GetResults() map[string]int {
	return fp.results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	processor := NewFileProcessor()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			processor.ProcessFile(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	processor.Wait()

	results := processor.GetResults()
	fmt.Println("File processing results:")
	for file, lines := range results {
		fmt.Printf("%s: %d lines\n", file, lines)
	}
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
	results  []string
	wg       sync.WaitGroup
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		results: make([]string, 0),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	fp.mu.Lock()
	fp.results = append(fp.results, fmt.Sprintf("%s: %d lines", filepath.Base(path), lineCount))
	fp.mu.Unlock()

	return scanner.Err()
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fp.wg.Add(1)
		go func(filename string) {
			defer fp.wg.Done()
			fullPath := filepath.Join(dirPath, filename)
			if err := fp.ProcessFile(fullPath); err != nil {
				fmt.Printf("Error processing %s: %v\n", filename, err)
			}
		}(file.Name())
	}

	fp.wg.Wait()
	return nil
}

func (fp *FileProcessor) GetResults() []string {
	return fp.results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	processor := NewFileProcessor()
	if err := processor.ProcessDirectory(os.Args[1]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processing results:")
	for _, result := range processor.GetResults() {
		fmt.Println(result)
	}
}package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type Data struct {
    Name  string `json:"name"`
    Value int    `json:"value"`
}

func readJSON(filename string) (*Data, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var data Data
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&data)
    if err != nil {
        return nil, err
    }
    return &data, nil
}

func writeJSON(filename string, data *Data) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(data)
}

func main() {
    data := &Data{Name: "test", Value: 42}
    err := writeJSON("output.json", data)
    if err != nil {
        fmt.Println("Error writing file:", err)
        return
    }

    readData, err := readJSON("output.json")
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }
    fmt.Printf("Read data: %+v\n", readData)
}