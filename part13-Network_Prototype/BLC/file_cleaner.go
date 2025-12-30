
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
}