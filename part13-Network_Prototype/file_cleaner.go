package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	tempDir := os.TempDir()
	cutoffTime := time.Now().AddDate(0, 0, -7)

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing old file: %s\n", path)
			os.RemoveAll(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error cleaning temp directory: %v\n", err)
	}
}