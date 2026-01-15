package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const tempFileAgeThreshold = 24 * time.Hour

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_cleaner <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	err := cleanTempFiles(dir)
	if err != nil {
		fmt.Printf("Error cleaning directory: %v\n", err)
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

		if isTempFile(info.Name()) && isOldFile(info.ModTime()) {
			fmt.Printf("Removing: %s\n", path)
			return os.Remove(path)
		}

		return nil
	})
}

func isTempFile(filename string) bool {
	return len(filename) > 4 && filename[len(filename)-4:] == ".tmp"
}

func isOldFile(modTime time.Time) bool {
	return time.Since(modTime) > tempFileAgeThreshold
}