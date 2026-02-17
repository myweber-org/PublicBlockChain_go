
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type LogRotator struct {
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
}

func NewLogRotator(filePath string, maxSizeMB int) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &LogRotator{
        filePath:    filePath,
        maxSize:     maxSize,
        currentSize: stat.Size(),
        file:        file,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    lr.file.Close()

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

    if err := os.Rename(lr.filePath, backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.file = file
    lr.currentSize = 0

    go lr.compressOldLogs()
    return nil
}

func (lr *LogRotator) compressOldLogs() {
    dir := filepath.Dir(lr.filePath)
    pattern := filepath.Base(lr.filePath) + ".*"

    matches, err := filepath.Glob(filepath.Join(dir, pattern))
    if err != nil {
        return
    }

    for _, match := range matches {
        if filepath.Ext(match) != ".gz" {
            lr.compressFile(match)
        }
    }
}

func (lr *LogRotator) compressFile(path string) {
    // Compression implementation would go here
    // For simplicity, we'll just log the intention
    fmt.Printf("Would compress: %s\n", path)
}

func (lr *LogRotator) Close() error {
    return lr.file.Close()
}

func main() {
    rotator, err := NewLogRotator("app.log", 10) // 10MB max size
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }
}