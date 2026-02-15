package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_temp"
	daysToKeep   = 7
	fileExt      = ".tmp"
)

func main() {
	err := cleanOldFiles(tempDir, daysToKeep, fileExt)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string, daysOld int, extension string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -daysOld)
	removedCount := 0

	for _, file := range files {
		if filepath.Ext(file.Name()) != extension {
			continue
		}

		if file.ModTime().Before(cutoffTime) {
			fullPath := filepath.Join(dirPath, file.Name())
			err := os.Remove(fullPath)
			if err != nil {
				fmt.Printf("Warning: Could not remove %s: %v\n", fullPath, err)
				continue
			}
			removedCount++
			fmt.Printf("Removed: %s\n", fullPath)
		}
	}

	fmt.Printf("Total files removed: %d\n", removedCount)
	return nil
}