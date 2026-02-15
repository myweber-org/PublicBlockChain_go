package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const retentionDays = 7

func main() {
	tempDir := os.TempDir()
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing old file: %s\n", path)
			os.Remove(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error cleaning temp directory: %v\n", err)
	}
}package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	daysToKeep = 7
	tempDir    = "/tmp/app_temp"
)

func main() {
	err := cleanOldFiles(tempDir)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dir string) error {
	cutoffTime := time.Now().AddDate(0, 0, -daysToKeep)

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
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
			err = os.Remove(path)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", path, err)
			} else {
				fmt.Printf("Removed old file: %s\n", path)
			}
		}

		return nil
	})
}