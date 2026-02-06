
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

		fmt.Printf("Processed %s: %d lines\n", path, lineCount)
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
	fmt.Printf("\nTotal files processed: %d\n", len(results))
}package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "sync"
)

type FileProcessor struct {
    inputDir  string
    outputDir string
    wg        sync.WaitGroup
}

func NewFileProcessor(input, output string) *FileProcessor {
    return &FileProcessor{
        inputDir:  input,
        outputDir: output,
    }
}

func (fp *FileProcessor) ProcessFile(filename string) error {
    defer fp.wg.Done()

    inputPath := filepath.Join(fp.inputDir, filename)
    outputPath := filepath.Join(fp.outputDir, "processed_"+filename)

    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("cannot open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("cannot create output file: %w", err)
    }
    defer outputFile.Close()

    scanner := bufio.NewScanner(inputFile)
    writer := bufio.NewWriter(outputFile)

    for scanner.Scan() {
        line := scanner.Text()
        processedLine := transformLine(line)
        _, err := writer.WriteString(processedLine + "\n")
        if err != nil {
            return fmt.Errorf("write error: %w", err)
        }
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("scan error: %w", err)
    }

    return writer.Flush()
}

func transformLine(line string) string {
    return "PROCESSED: " + line
}

func (fp *FileProcessor) ProcessAll() []error {
    entries, err := os.ReadDir(fp.inputDir)
    if err != nil {
        return []error{fmt.Errorf("cannot read directory: %w", err)}
    }

    var errors []error
    errorChan := make(chan error, len(entries))

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        fp.wg.Add(1)
        go func(fname string) {
            if err := fp.ProcessFile(fname); err != nil {
                errorChan <- err
            }
        }(entry.Name())
    }

    fp.wg.Wait()
    close(errorChan)

    for err := range errorChan {
        errors = append(errors, err)
    }

    return errors
}

func main() {
    processor := NewFileProcessor("./input", "./output")
    errors := processor.ProcessAll()

    if len(errors) > 0 {
        fmt.Printf("Processing completed with %d errors:\n", len(errors))
        for _, err := range errors {
            fmt.Println(err)
        }
    } else {
        fmt.Println("All files processed successfully")
    }
}