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
}package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadJSONFile reads a JSON file and unmarshals it into the provided interface.
func ReadJSONFile(filename string, v interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	return nil
}

// WriteJSONFile marshals the provided data and writes it to a file.
func WriteJSONFile(filename string, v interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

func main() {
	// Example usage
	type Config struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	config := Config{Host: "localhost", Port: 8080}
	if err := WriteJSONFile("config.json", config); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}
	fmt.Println("File written successfully")

	var loadedConfig Config
	if err := ReadJSONFile("config.json", &loadedConfig); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	fmt.Printf("Loaded config: %+v\n", loadedConfig)
}