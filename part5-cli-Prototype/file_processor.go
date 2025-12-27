package main

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
    workers   int
}

func NewFileProcessor(input, output string, workers int) *FileProcessor {
    return &FileProcessor{
        inputDir:  input,
        outputDir: output,
        workers:   workers,
    }
}

func (fp *FileProcessor) Process() error {
    files, err := os.ReadDir(fp.inputDir)
    if err != nil {
        return fmt.Errorf("failed to read input directory: %w", err)
    }

    jobs := make(chan string, len(files))
    results := make(chan error, len(files))
    var wg sync.WaitGroup

    for w := 0; w < fp.workers; w++ {
        wg.Add(1)
        go fp.worker(jobs, results, &wg)
    }

    for _, file := range files {
        if !file.IsDir() {
            jobs <- file.Name()
        }
    }
    close(jobs)

    wg.Wait()
    close(results)

    for err := range results {
        if err != nil {
            return err
        }
    }
    return nil
}

func (fp *FileProcessor) worker(jobs <-chan string, results chan<- error, wg *sync.WaitGroup) {
    defer wg.Done()
    for filename := range jobs {
        err := fp.processFile(filename)
        results <- err
    }
}

func (fp *FileProcessor) processFile(filename string) error {
    inputPath := filepath.Join(fp.inputDir, filename)
    outputPath := filepath.Join(fp.outputDir, "processed_"+filename)

    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    scanner := bufio.NewScanner(inputFile)
    writer := bufio.NewWriter(outputFile)

    for scanner.Scan() {
        line := scanner.Text()
        processed := transformLine(line)
        if _, err := writer.WriteString(processed + "\n"); err != nil {
            return fmt.Errorf("failed to write line: %w", err)
        }
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading file: %w", err)
    }

    if err := writer.Flush(); err != nil {
        return fmt.Errorf("failed to flush writer: %w", err)
    }

    return nil
}

func transformLine(line string) string {
    var result []rune
    for _, r := range line {
        if r >= 'a' && r <= 'z' {
            result = append(result, r-32)
        } else if r >= 'A' && r <= 'Z' {
            result = append(result, r+32)
        } else {
            result = append(result, r)
        }
    }
    return string(result)
}

func main() {
    processor := NewFileProcessor("./input", "./output", 4)
    if err := processor.Process(); err != nil {
        fmt.Printf("Processing failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("File processing completed successfully")
}