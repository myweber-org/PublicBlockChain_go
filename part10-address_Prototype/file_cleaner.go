package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	dir := "./tmp"
	maxAge := time.Hour * 24 * 7

	err := cleanOldFiles(dir, maxAge)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		return
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dir string, maxAge time.Duration) error {
	cutoff := time.Now().Add(-maxAge)

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoff) {
			fmt.Printf("Removing old file: %s\n", path)
			return os.Remove(path)
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

const (
	tempDir      = "/tmp/app_cache"
	retentionDays = 7
)

func main() {
	err := cleanOldFiles(tempDir)
	if err != nil {
		fmt.Printf("Cleanup failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string) error {
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

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
			err = os.Remove(path)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", path, err)
			} else {
				fmt.Printf("Removed: %s\n", path)
			}
		}
		return nil
	})
}