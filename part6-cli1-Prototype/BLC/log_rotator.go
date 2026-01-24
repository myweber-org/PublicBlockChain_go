
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    currentFile   *os.File
    currentSize   int64
    basePath      string
    currentIndex  int
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }

    if err := lr.compressOldLogs(); err != nil {
        return err
    }

    lr.currentIndex++
    return lr.openCurrentFile()
}

func (lr *LogRotator) openCurrentFile() error {
    filename := fmt.Sprintf("%s.%d", lr.basePath, lr.currentIndex)
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = stat.Size()
    return nil
}

func (lr *LogRotator) compressOldLogs() error {
    for i := lr.currentIndex - maxBackups; i >= 0; i-- {
        logPath := fmt.Sprintf("%s.%d", lr.basePath, i)
        compressedPath := logPath + ".gz"

        if _, err := os.Stat(compressedPath); err == nil {
            continue
        }

        if err := compressFile(logPath, compressedPath); err != nil {
            return err
        }

        os.Remove(logPath)
    }
    return nil
}

func compressFile(src, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func findLatestLogIndex(basePath string) int {
    pattern := basePath + ".*"
    matches, _ := filepath.Glob(pattern)

    maxIndex := 0
    for _, match := range matches {
        if strings.HasSuffix(match, ".gz") {
            continue
        }

        parts := strings.Split(match, ".")
        if len(parts) < 2 {
            continue
        }

        idx, err := strconv.Atoi(parts[len(parts)-1])
        if err == nil && idx > maxIndex {
            maxIndex = idx
        }
    }
    return maxIndex
}

func main() {
    logPath := "application.log"
    rotator, err := NewLogRotator(logPath)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: This is a sample log message\n",
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation completed successfully")
}