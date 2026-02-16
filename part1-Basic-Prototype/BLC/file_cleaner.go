package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	LogDir       = "./logs"
	MaxAgeDays   = 30
	ScanInterval = 24 * time.Hour
)

func main() {
	ticker := time.NewTicker(ScanInterval)
	defer ticker.Stop()

	fmt.Println("Starting log cleanup service...")
	cleanupOldLogs()

	for range ticker.C {
		cleanupOldLogs()
	}
}

func cleanupOldLogs() {
	cutoffTime := time.Now().AddDate(0, 0, -MaxAgeDays)

	err := filepath.WalkDir(LogDir, func(path string, d fs.DirEntry, err error) error {
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
			if err := os.Remove(path); err != nil {
				fmt.Printf("Failed to remove %s: %v\n", path, err)
			} else {
				fmt.Printf("Removed old log file: %s\n", path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
}package main

import (
    "os"
    "path/filepath"
    "time"
)

func main() {
    dir := "/tmp"
    cutoff := time.Now().AddDate(0, 0, -7)

    filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
        if info.IsDir() {
            return nil
        }
        if info.ModTime().Before(cutoff) {
            os.Remove(path)
        }
        return nil
    })
}