package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	tempFilePrefix = "temp_"
	maxAgeHours    = 24
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_cleaner <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	err := cleanTempFiles(dir)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Cleanup completed successfully")
}

func cleanTempFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !isTempFile(info.Name()) {
			return nil
		}

		if isFileOld(info.ModTime()) {
			return os.Remove(path)
		}

		return nil
	})
}

func isTempFile(filename string) bool {
	return len(filename) > len(tempFilePrefix) && filename[:len(tempFilePrefix)] == tempFilePrefix
}

func isFileOld(modTime time.Time) bool {
	age := time.Since(modTime)
	return age > maxAgeHours*time.Hour
}