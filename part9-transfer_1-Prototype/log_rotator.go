package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackupCount = 5
    logFileName   = "app.log"
)

type LogRotator struct {
    currentSize int64
    basePath    string
}

func NewLogRotator(basePath string) *LogRotator {
    return &LogRotator{
        basePath: basePath,
    }
}

func (lr *LogRotator) Write(p []byte) (n int, err error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    file, err := os.OpenFile(filepath.Join(lr.basePath, logFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    n, err = file.Write(p)
    lr.currentSize += int64(n)
    return n, err
}

func (lr *LogRotator) rotate() error {
    currentPath := filepath.Join(lr.basePath, logFileName)
    if _, err := os.Stat(currentPath); os.IsNotExist(err) {
        return nil
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := filepath.Join(lr.basePath, fmt.Sprintf("%s.%s", logFileName, timestamp))

    if err := os.Rename(currentPath, backupPath); err != nil {
        return err
    }

    lr.currentSize = 0
    lr.cleanupOldBackups()
    return nil
}

func (lr *LogRotator) cleanupOldBackups() {
    files, err := filepath.Glob(filepath.Join(lr.basePath, logFileName+".*"))
    if err != nil {
        return
    }

    sort.Sort(sort.Reverse(sort.StringSlice(files)))

    for i, file := range files {
        if i >= maxBackupCount {
            os.Remove(file)
        }
    }
}

func (lr *LogRotator) loadCurrentSize() error {
    fileInfo, err := os.Stat(filepath.Join(lr.basePath, logFileName))
    if os.IsNotExist(err) {
        lr.currentSize = 0
        return nil
    }
    if err != nil {
        return err
    }
    lr.currentSize = fileInfo.Size()
    return nil
}

func main() {
    rotator := NewLogRotator(".")
    if err := rotator.loadCurrentSize(); err != nil {
        fmt.Printf("Failed to load current log size: %v\n", err)
        return
    }

    testMessage := fmt.Sprintf("Test log entry at %s\n", time.Now().Format(time.RFC3339))
    if _, err := rotator.Write([]byte(testMessage)); err != nil {
        fmt.Printf("Failed to write log: %v\n", err)
    } else {
        fmt.Println("Log entry written successfully")
    }
}