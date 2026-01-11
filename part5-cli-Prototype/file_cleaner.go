
package main

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
	err := cleanOldFiles(tempDir, maxAgeHours)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string, maxAgeHours int) error {
	cutoffTime := time.Now().Add(-time.Duration(maxAgeHours) * time.Hour)

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing old file: %s (modified: %v)\n", path, info.ModTime())
			return os.Remove(path)
		}

		return nil
	})
}