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
}