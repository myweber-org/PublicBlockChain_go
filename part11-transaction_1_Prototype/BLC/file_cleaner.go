package main

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