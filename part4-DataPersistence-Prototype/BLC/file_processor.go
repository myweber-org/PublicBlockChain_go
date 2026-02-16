
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
}