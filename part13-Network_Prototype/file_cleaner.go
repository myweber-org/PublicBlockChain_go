package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const retentionDays = 7

func main() {
	tempDir := os.TempDir()
	fmt.Printf("Cleaning temporary files in: %s\n", tempDir)

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	var removedCount int
	var totalSize int64

	err := filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			size := info.Size()
			if err := os.Remove(path); err == nil {
				removedCount++
				totalSize += size
				fmt.Printf("Removed: %s (size: %d bytes)\n", filepath.Base(path), size)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error during cleanup: %v\n", err)
	}

	fmt.Printf("Cleanup completed. Removed %d files, freed %d bytes.\n", removedCount, totalSize)
}package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_cache"
	maxAgeHours  = 168
)

func main() {
	err := cleanOldFiles(tempDir, maxAgeHours)
	if err != nil {
		fmt.Printf("Cleanup failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string, maxAgeHours int) error {
	cutoffTime := time.Now().Add(-time.Duration(maxAgeHours) * time.Hour)

	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing old file: %s (modified: %v)\n", path, info.ModTime())
			return os.Remove(path)
		}

		return nil
	})
}