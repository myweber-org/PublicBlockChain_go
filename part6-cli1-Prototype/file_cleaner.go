package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const tempDir = "/tmp/myapp"
const retentionDays = 7

func main() {
	err := cleanOldFiles(tempDir, retentionDays)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string, days int) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	cutoffTime := time.Now().AddDate(0, 0, -days)
	removedCount := 0

	for _, file := range files {
		if file.ModTime().Before(cutoffTime) {
			fullPath := filepath.Join(dirPath, file.Name())
			err := os.Remove(fullPath)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", fullPath, err)
			} else {
				removedCount++
			}
		}
	}

	fmt.Printf("Removed %d files older than %d days\n", removedCount, days)
	return nil
}