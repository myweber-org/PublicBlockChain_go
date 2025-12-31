
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
)

type LogRotator struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    sequence    int
}

func NewLogRotator(basePath string, maxSize int64) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
        maxSize:  maxSize,
        sequence: 0,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.file != nil {
        lr.file.Close()
    }

    file, err := os.OpenFile(lr.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.file = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize && lr.currentSize > 0 {
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
    if lr.file != nil {
        lr.file.Close()
        lr.file = nil
    }

    rotatedPath := lr.basePath + "." + strconv.Itoa(lr.sequence)
    for {
        if _, err := os.Stat(rotatedPath); os.IsNotExist(err) {
            break
        }
        lr.sequence++
        rotatedPath = lr.basePath + "." + strconv.Itoa(lr.sequence)
    }

    if err := os.Rename(lr.basePath, rotatedPath); err != nil {
        return err
    }

    go lr.compressOldLog(rotatedPath)

    return lr.openCurrentFile()
}

func (lr *LogRotator) compressOldLog(path string) {
    originalFile, err := os.Open(path)
    if err != nil {
        return
    }
    defer originalFile.Close()

    compressedPath := path + ".gz"
    compressedFile, err := os.Create(compressedPath)
    if err != nil {
        return
    }
    defer compressedFile.Close()

    gzWriter := gzip.NewWriter(compressedFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, originalFile); err != nil {
        return
    }

    os.Remove(path)
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.file != nil {
        return lr.file.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log", 1024*1024)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 10000; i++ {
        logEntry := fmt.Sprintf("Log entry number %d: Some sample log data here\n", i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
    }

    fmt.Println("Log rotation test completed")
}