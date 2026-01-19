package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileStats struct {
	Path    string
	Size    int64
	Lines   int
	Words   int
	Chars   int
}

func processFile(path string, wg *sync.WaitGroup, results chan<- FileStats) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", path, err)
		return
	}
	defer file.Close()

	stats := FileStats{Path: path}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info for %s: %v\n", path, err)
		return
	}
	stats.Size = fileInfo.Size()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stats.Lines++
		stats.Chars += len(line)
		
		wordScanner := bufio.NewScanner(bufio.NewReader(filepath.NewReader(line)))
		wordScanner.Split(bufio.ScanWords)
		for wordScanner.Scan() {
			stats.Words++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file %s: %v\n", path, err)
		return
	}

	results <- stats
}

func collectResults(results <-chan FileStats, totalFiles *int, totalStats *FileStats) {
	for stats := range results {
		*totalFiles++
		totalStats.Size += stats.Size
		totalStats.Lines += stats.Lines
		totalStats.Words += stats.Words
		totalStats.Chars += stats.Chars

		fmt.Printf("Processed: %s\n", stats.Path)
		fmt.Printf("  Size: %d bytes, Lines: %d, Words: %d, Characters: %d\n\n",
			stats.Size, stats.Lines, stats.Words, stats.Chars)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_processor.go <file1> [file2] ...")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	results := make(chan FileStats, len(os.Args)-1)
	totalFiles := 0
	totalStats := FileStats{}

	go collectResults(results, &totalFiles, &totalStats)

	for _, filePath := range os.Args[1:] {
		wg.Add(1)
		go processFile(filePath, &wg, results)
	}

	wg.Wait()
	close(results)

	fmt.Printf("\n=== SUMMARY ===\n")
	fmt.Printf("Total files processed: %d\n", totalFiles)
	fmt.Printf("Combined size: %d bytes\n", totalStats.Size)
	fmt.Printf("Total lines: %d\n", totalStats.Lines)
	fmt.Printf("Total words: %d\n", totalStats.Words)
	fmt.Printf("Total characters: %d\n", totalStats.Chars)
}