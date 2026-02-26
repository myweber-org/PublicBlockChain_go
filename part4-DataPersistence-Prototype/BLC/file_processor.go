
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

	results := processor.GetResults()
	fmt.Println("Processing results:")
	for _, result := range results {
		fmt.Println(result)
	}
}package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type FileProcessor struct {
	inputPath  string
	outputPath string
	mu         sync.Mutex
}

func NewFileProcessor(input, output string) *FileProcessor {
	return &FileProcessor{
		inputPath:  input,
		outputPath: output,
	}
}

func (fp *FileProcessor) ProcessLines() error {
	file, err := os.Open(fp.inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	outputFile, err := os.Create(fp.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outputFile)
	var wg sync.WaitGroup
	lineChan := make(chan string, 100)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lineChan {
				processed := fp.transformLine(line)
				fp.mu.Lock()
				writer.WriteString(processed + "\n")
				fp.mu.Unlock()
			}
		}()
	}

	for scanner.Scan() {
		lineChan <- scanner.Text()
	}
	close(lineChan)

	wg.Wait()
	writer.Flush()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	return nil
}

func (fp *FileProcessor) transformLine(line string) string {
	return fmt.Sprintf("Processed: %s", line)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: file_processor <input_file> <output_file>")
		os.Exit(1)
	}

	processor := NewFileProcessor(os.Args[1], os.Args[2])
	if err := processor.ProcessLines(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File processing completed successfully")
}package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type FileProcessor struct {
	FilePath string
}

func NewFileProcessor(path string) *FileProcessor {
	return &FileProcessor{FilePath: path}
}

func (fp *FileProcessor) Read() (string, error) {
	data, err := ioutil.ReadFile(fp.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

func (fp *FileProcessor) Write(content string) error {
	err := ioutil.WriteFile(fp.FilePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func (fp *FileProcessor) Append(content string) error {
	file, err := os.OpenFile(fp.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for append: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to append to file: %w", err)
	}
	return nil
}

func main() {
	processor := NewFileProcessor("example.txt")
	
	err := processor.Write("Initial content\n")
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}
	
	err = processor.Append("Appended content\n")
	if err != nil {
		fmt.Printf("Append error: %v\n", err)
		return
	}
	
	content, err := processor.Read()
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
		return
	}
	
	fmt.Printf("File content:\n%s", content)
}package main

import (
	"fmt"
	"sync"
	"time"
)

type FileProcessor struct {
	mu       sync.Mutex
	queue    []string
	workers  int
	done     chan bool
	wg       sync.WaitGroup
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		queue:   make([]string, 0),
		workers: workers,
		done:    make(chan bool),
	}
}

func (fp *FileProcessor) AddFile(filePath string) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	fp.queue = append(fp.queue, filePath)
}

func (fp *FileProcessor) processFile(filePath string) {
	defer fp.wg.Done()
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Processed: %s\n", filePath)
}

func (fp *FileProcessor) worker(id int) {
	for {
		select {
		case <-fp.done:
			return
		default:
			fp.mu.Lock()
			if len(fp.queue) > 0 {
				file := fp.queue[0]
				fp.queue = fp.queue[1:]
				fp.mu.Unlock()
				fp.wg.Add(1)
				go fp.processFile(file)
			} else {
				fp.mu.Unlock()
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
}

func (fp *FileProcessor) Start() {
	for i := 0; i < fp.workers; i++ {
		go fp.worker(i)
	}
}

func (fp *FileProcessor) Stop() {
	close(fp.done)
	fp.wg.Wait()
}

func main() {
	processor := NewFileProcessor(3)
	processor.Start()

	for i := 1; i <= 10; i++ {
		processor.AddFile(fmt.Sprintf("file%d.txt", i))
	}

	time.Sleep(2 * time.Second)
	processor.Stop()
	fmt.Println("All files processed")
}