
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_temp"
	retentionDays = 7
)

func main() {
	if err := cleanOldFiles(); err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles() error {
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	return filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if info.ModTime().Before(cutoffTime) {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove %s: %w", path, err)
			}
			fmt.Printf("Removed: %s\n", path)
		}
		return nil
	})
}package main

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