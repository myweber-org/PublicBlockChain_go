package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    maxFiles    int
    currentSize int64
    file        *os.File
}

func NewRotator(filePath string, maxSize int64, maxFiles int) (*Rotator, error) {
    r := &Rotator{
        filePath: filePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
        return nil, err
    }

    if err := r.openCurrentFile(); err != nil {
        return nil, err
    }

    go r.timeBasedRotation()
    return r, nil
}

func (r *Rotator) openCurrentFile() error {
    file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    r.file = file
    r.currentSize = stat.Size()
    return nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize+int64(len(p)) > r.maxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.file.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *Rotator) rotate() error {
    if r.file != nil {
        r.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)

    if err := os.Rename(r.filePath, backupPath); err != nil {
        return err
    }

    if err := r.cleanupOldFiles(); err != nil {
        fmt.Printf("Cleanup error: %v\n", err)
    }

    return r.openCurrentFile()
}

func (r *Rotator) cleanupOldFiles() error {
    pattern := fmt.Sprintf("%s.*", r.filePath)
    files, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(files) <= r.maxFiles {
        return nil
    }

    for i := 0; i < len(files)-r.maxFiles; i++ {
        if err := os.Remove(files[i]); err != nil {
            return err
        }
    }
    return nil
}

func (r *Rotator) timeBasedRotation() {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        r.mu.Lock()
        if r.currentSize > 0 {
            if err := r.rotate(); err != nil {
                fmt.Printf("Time-based rotation failed: %v\n", err)
            }
        }
        r.mu.Unlock()
    }
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.file != nil {
        return r.file.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("/var/log/myapp/app.log", 10*1024*1024, 5)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(100 * time.Millisecond)
    }
}