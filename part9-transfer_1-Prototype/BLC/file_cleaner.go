package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	tempDir := os.TempDir()
	cutoff := time.Now().AddDate(0, 0, -7)
	var removedCount int

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err == nil {
				removedCount++
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	fmt.Printf("Cleaned %d temporary files older than 7 days\n", removedCount)
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_temp"
	maxAgeHours  = 168 // 7 days
)

func main() {
	err := cleanOldFiles(tempDir)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string) error {
	cutoffTime := time.Now().Add(-time.Hour * maxAgeHours)

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			err := os.Remove(path)
			if err != nil {
				return fmt.Errorf("failed to remove %s: %w", path, err)
			}
			fmt.Printf("Removed: %s\n", path)
		}
		return nil
	})
}