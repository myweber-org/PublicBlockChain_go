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
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		stats.Lines++
		stats.Chars += len(line) + 1
		stats.Words += countWords(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file %s: %v\n", path, err)
		return
	}

	if fileInfo, err := file.Stat(); err == nil {
		stats.Size = fileInfo.Size()
	}

	results <- stats
}

func countWords(line string) int {
	inWord := false
	count := 0

	for _, ch := range line {
		if ch == ' ' || ch == '\t' || ch == '\n' {
			inWord = false
		} else if !inWord {
			inWord = true
			count++
		}
	}

	return count
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	root := os.Args[1]
	var wg sync.WaitGroup
	results := make(chan FileStats, 100)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			wg.Add(1)
			go processFile(path, &wg, results)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalFiles := 0
	totalLines := 0
	totalWords := 0
	totalChars := 0

	for stats := range results {
		totalFiles++
		totalLines += stats.Lines
		totalWords += stats.Words
		totalChars += stats.Chars

		fmt.Printf("File: %s\n", stats.Path)
		fmt.Printf("  Size: %d bytes\n", stats.Size)
		fmt.Printf("  Lines: %d\n", stats.Lines)
		fmt.Printf("  Words: %d\n", stats.Words)
		fmt.Printf("  Characters: %d\n\n", stats.Chars)
	}

	fmt.Printf("Summary:\n")
	fmt.Printf("  Total files processed: %d\n", totalFiles)
	fmt.Printf("  Total lines: %d\n", totalLines)
	fmt.Printf("  Total words: %d\n", totalWords)
	fmt.Printf("  Total characters: %d\n", totalChars)
}