package main

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
	err := cleanOldFiles(tempDir, retentionDays)
	if err != nil {
		fmt.Printf("Cleanup failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string, days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)

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
			fmt.Printf("Removing: %s (modified: %v)\n", path, info.ModTime())
			return os.Remove(path)
		}
		return nil
	})
}
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Cleaner struct {
	Directory   string
	OlderThan   time.Duration
	Extensions  []string
	DryRun      bool
}

func NewCleaner(dir string, olderThan time.Duration, exts []string, dryRun bool) *Cleaner {
	return &Cleaner{
		Directory:  dir,
		OlderThan:  olderThan,
		Extensions: exts,
		DryRun:     dryRun,
	}
}

func (c *Cleaner) ShouldDelete(filePath string, info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	cutoffTime := time.Now().Add(-c.OlderThan)
	if info.ModTime().After(cutoffTime) {
		return false
	}

	if len(c.Extensions) > 0 {
		ext := filepath.Ext(filePath)
		matched := false
		for _, allowedExt := range c.Extensions {
			if ext == allowedExt {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}

func (c *Cleaner) Clean() (int, error) {
	var count int

	err := filepath.Walk(c.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if c.ShouldDelete(path, info) {
			if c.DryRun {
				fmt.Printf("[DRY RUN] Would delete: %s (modified: %v)\n", path, info.ModTime())
			} else {
				if err := os.Remove(path); err != nil {
					return fmt.Errorf("failed to delete %s: %w", path, err)
				}
				fmt.Printf("Deleted: %s\n", path)
			}
			count++
		}
		return nil
	})

	return count, err
}

func main() {
	cleaner := NewCleaner(
		"/tmp",
		24*time.Hour,
		[]string{".log", ".tmp", ".cache"},
		true,
	)

	count, err := cleaner.Clean()
	if err != nil {
		fmt.Printf("Error during cleanup: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processed %d files\n", count)
}