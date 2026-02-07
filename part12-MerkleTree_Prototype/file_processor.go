
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FileInfo struct {
	Path string
	Size int64
	Hash string
	Err  error
}

func processFile(path string, results chan<- FileInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	info := FileInfo{Path: path}

	file, err := os.Open(path)
	if err != nil {
		info.Err = err
		results <- info
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		info.Err = err
		results <- info
		return
	}
	info.Size = stat.Size()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		info.Err = err
		results <- info
		return
	}

	info.Hash = hex.EncodeToString(hasher.Sum(nil))
	results <- info
}

func validateFiles(dir string) ([]FileInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	results := make(chan FileInfo, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		wg.Add(1)
		go processFile(filepath.Join(dir, entry.Name()), results, &wg)
	}

	wg.Wait()
	close(results)

	var processed []FileInfo
	for info := range results {
		processed = append(processed, info)
	}

	return processed, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	files, err := validateFiles(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processed %d files:\n", len(files))
	for _, file := range files {
		if file.Err != nil {
			fmt.Printf("ERROR %s: %v\n", file.Path, file.Err)
		} else {
			fmt.Printf("OK %s: size=%d hash=%s\n", file.Path, file.Size, file.Hash)
		}
	}
}