package main

import (
    "os"
    "path/filepath"
    "time"
)

func main() {
    tempDir := os.TempDir()
    cutoff := time.Now().AddDate(0, 0, -7)

    filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
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
}package main

import (
    "os"
    "path/filepath"
    "time"
)

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
            return os.Remove(path)
        }
        return nil
    })
}

func main() {
    tempDir := os.TempDir()
    err := cleanOldFiles(tempDir, 7*24*time.Hour)
    if err != nil {
        panic(err)
    }
}