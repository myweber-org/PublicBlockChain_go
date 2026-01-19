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
}package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	retentionDays = 7
	tempDir       = "/tmp/myapp"
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
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
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
			fmt.Printf("Removing old file: %s (modified: %v)\n", path, info.ModTime())
			return os.Remove(path)
		}

		return nil
	})
}