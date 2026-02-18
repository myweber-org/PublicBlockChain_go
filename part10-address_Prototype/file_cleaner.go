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
}