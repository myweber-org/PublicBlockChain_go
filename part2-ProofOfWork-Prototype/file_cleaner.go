package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func main() {
	tempDir := os.TempDir()
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		fmt.Printf("Error reading temp directory: %v\n", err)
		return
	}

	cutoffTime := time.Now().Add(-7 * 24 * time.Hour)
	removedCount := 0

	for _, file := range files {
		if file.ModTime().Before(cutoffTime) {
			filePath := filepath.Join(tempDir, file.Name())
			err := os.RemoveAll(filePath)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", filePath, err)
			} else {
				removedCount++
			}
		}
	}

	fmt.Printf("Cleaned up %d old temporary files\n", removedCount)
}